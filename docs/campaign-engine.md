# Central de Campanhas / Outbound Multicanal — Especificação Técnica

> Status: **rascunho para revisão** · Não implementado · Não afeta nada em produção.
> Objetivo: definir o contrato (arquitetura, schema, motor, drivers, throttle) antes de escrever código.

---

## 1. Visão geral

Módulo de disparo outbound (marketing/reativação) de **alto nível** para o Vysor CRM, multicanal
(WhatsApp, Telegram, Email, Instagram, Facebook), com foco em **não tomar ban/bloqueio** através de
envio com cadência humana, warm-up de números e respeito às regras de cada canal.

### Princípios inegociáveis
1. **Serviço isolado e aditivo.** Nada do que roda hoje (agentes, TTS, handoff, fluxo de vídeo) é
   alterado. Se o motor cair, o CRM continua funcionando.
2. **O ritmo mora no motor.** n8n e evolution-api são apenas "braços" que enviam 1 mensagem por vez,
   quando o motor autoriza. Nunca se manda uma lista pra um braço "resolver".
3. **Núcleo agnóstico + driver por canal.** O motor não sabe o que é WhatsApp ou Email; ele fala com
   uma interface comum. Cada canal é um plugin com suas próprias regras.
4. **A inbox é a fonte de verdade do canal.** A campanha escolhe uma inbox; o `channel_type` dela
   decide o driver. Cadastrou número novo em Canais → aparece no seletor automaticamente.

---

## 2. Arquitetura

```
CRM/Rails (cria campanha + audiência)
      → Redis (fila + agenda por timestamp)   [DB dedicado: /6]
      → evo-campaign-engine (Go)  ── decide QUANDO e PRA QUEM ──┐
            │                                                    │
            ├─ WhatsApp texto ───────────→ evolution-api         │
            ├─ WhatsApp mídia ───────────→ n8n webhook (MinIO)   │
            └─ Telegram/Email/IG/FB ─────→ plumbing do CRM       │
                                                                 │
      ← webhooks de entrega ← evolution-api / CRM ───────────────┘
      → feedback loop (ban detection / auto-pause)
      ↕ personalização de texto (Python / evo-processor)
```

### Componentes e responsabilidade

| Camada | Serviço | O quê |
|---|---|---|
| Cérebro | **evo-campaign-engine** (Go, novo) | agenda, jitter, teto por número, warm-up, rotação, retry, ban detection |
| Criação/UI | evo-crm (Rails) + frontend (React) | CRUD campanha, seletor de canal, segmentação, relatórios |
| Personalização | evo-processor (Python) | variação de texto por contato (spintax/LLM) |
| Atuador WA texto | evolution-api | envio + presença "digitando" |
| Atuador WA mídia | **n8n** (fluxo atual) | `{telefone, video_url, media, instance}` |
| Atuador demais canais | plumbing do CRM | reusa entrega existente do Chatwoot |
| Estado/agenda | Redis (DB /6) + Postgres | filas, buckets, locks + estado durável |

**Por que serviço novo (e não dentro do evo-bot-runtime):** bot-runtime é tempo-real (chat ao vivo);
o motor é cron/fila/pausável. Perfis de carga e falha diferentes → blast radius isolado, deploy e
pause independentes. Reusa a base Go do bot-runtime (libs, padrões, Redis).

---

## 3. Modelo de dados (Postgres)

Tabelas próprias do serviço, prefixo `cmp_` (mesma instância `evo_community`, sem tocar no schema atual).

```sql
-- Perfis de ritmo reutilizáveis, por tipo de canal
cmp_throttle_profiles(
  id, name, channel_type,
  daily_cap_start int, daily_cap_max int,        -- warm-up
  warmup_multiplier numeric, warmup_step_days int,
  hourly_cap int,
  min_delay_sec int, max_delay_sec int,          -- jitter entre envios
  coffee_break_every_n int, coffee_break_min_sec int, coffee_break_max_sec int,
  quiet_hours_start time, quiet_hours_end time, respect_timezone bool,
  created_at, updated_at
)

-- Número/handle de envio, ligado a uma inbox do CRM
cmp_sender_instances(
  id, inbox_id,               -- referencia a inbox (canal) do CRM
  channel_type,               -- denormalizado da inbox
  purpose,                    -- 'conversation' | 'campaign' (só rótulo)
  identifier,                 -- telefone/email/handle
  warmup_started_at, warmup_stage_day int,
  health_score int,           -- 0..100
  status,                     -- 'active' | 'paused' | 'banned'
  last_sent_at,
  created_at, updated_at
)

cmp_campaigns(
  id, account_id, name,
  trigger_type,               -- 'manual' | 'on_label' | 'scheduled'
  trigger_config jsonb,       -- label_id, cron, wait_hours
  throttle_profile_id,
  status,                     -- draft|scheduled|running|paused|completed|cancelled
  created_by, created_at, updated_at, started_at, finished_at
)

-- Pool de envio: 1..N inboxes escolhidas no seletor (multi-check) → rotação
cmp_campaign_channels(
  id, campaign_id, inbox_id
)  -- UNIQUE(campaign_id, inbox_id)

cmp_campaign_audience(
  id, campaign_id,
  mode,                       -- 'static' | 'dynamic'
  segment_filter jsonb,       -- para dynamic (tag/última compra/etc.)
  total_contacts int, created_at
)

cmp_audience_members(         -- snapshot resolvido (modo static)
  id, campaign_id, contact_id,
  recipient,                  -- telefone/email/handle denormalizado
  timezone, state,            -- pending|queued|done|skipped
)

cmp_message_variants(         -- variações de conteúdo (spintax / A-B)
  id, campaign_id,
  subject,                    -- email
  body, media_url, media_type,
  weight int, created_at
)

cmp_send_jobs(                -- unidade de trabalho: 1 por contato por campanha
  id, campaign_id, contact_id,
  sender_instance_id,         -- atribuído no disparo
  recipient, variant_id,
  state,                      -- queued|scheduled|sending|sent|delivered|read|failed|skipped|opted_out
  scheduled_at, sent_at, delivered_at, read_at,
  attempts int, last_error, created_at, updated_at
)

cmp_suppression(             -- opt-out / bounce / manual (por canal)
  id, account_id, contact_id, channel_type, reason, created_at
)  -- UNIQUE(account_id, contact_id, channel_type)

cmp_delivery_events(         -- feed cru dos webhooks (ban detection + métricas)
  id, send_job_id, status, provider_status, ts
)
```

Índices-chave: `cmp_send_jobs(campaign_id, state)`, `cmp_send_jobs(scheduled_at)`,
`cmp_suppression(account_id, contact_id, channel_type)` UNIQUE.

### Estruturas no Redis (DB /6)
- `cmp:schedule` — ZSET, score = timestamp agendado, member = job_id (jobs vencidos)
- `cmp:bucket:{instance}:day` / `:hour` — contadores com TTL (teto diário/horário)
- `cmp:lock:instance:{instance}` — serializa envio por número (cadência humana)
- `cmp:next_slot:{instance}` — timestamp do próximo envio permitido daquele número

---

## 4. Máquina de estados do `send_job`

```
queued → scheduled → sending → sent → delivered → read
                         ↘ failed → (retry: scheduled | desiste)
queued/scheduled → skipped        (suppression / opted_out)
```

---

## 5. Motor de throttle (o coração anti-ban)

### Warm-up (teto do número no dia D desde o início do aquecimento)
```
cap(D) = min(daily_cap_max,
             floor(daily_cap_start * warmup_multiplier ^ floor(D / warmup_step_days)))
# ex: start=20, mult=1.3, step=1  →  d0=20, d1=26, d2=33, ... até daily_cap_max
```

### Jitter (intervalo entre envios do MESMO número)
```
delay = random(min_delay_sec, max_delay_sec)          # ex: 45..120s
if enviados_desde_pausa >= coffee_break_every_n:
    delay += random(coffee_break_min_sec, coffee_break_max_sec)   # "pausa de café"
    zera enviados_desde_pausa
next_slot[instance] = now + delay
```

### Quiet hours / timezone
Se `now` (no fuso do contato) estiver fora da janela comercial → reagenda pra próxima abertura.

### Loop de disparo (scheduler + workers)
```
# Scheduler: promove jobs vencidos
a cada tick:
  due = ZRANGEBYSCORE cmp:schedule -inf now
  para cada job em due: envia pra fila de trabalho (stream); ZREM

# Worker:
job = lê da stream
inst = escolhe_numero(pool_do_canal)      # menos carregado e saudável
se inst == None:                 reagenda(job, +60s); ack; continua
se now < next_slot[inst]:        reagenda(job, next_slot[inst]); ack; continua   # cadência
se not bucket_allow(inst):       reagenda(job, próx_janela_bucket); ack; continua
se fora_horario(job.tz):         reagenda(job, próx_janela); ack; continua
se suppressed(recipient):        marca skipped(opted_out); ack; continua

driver = driver_do(channel_type)
texto  = personaliza(escolhe_variante(camp), job.contact)   # spintax local / LLM Python
driver.presence(inst, recipient, tempo_digitando(texto))    # só WhatsApp
res = driver.send(inst, recipient, conteudo)                # WA: evolution/n8n; resto: CRM
registra delivery_event(job, res)
se res.ok: marca sent; bucket_consume(inst); next_slot[inst] = now + jitter()
senão:     attempts++; se attempts<max: reagenda backoff; senão marca failed
ack
```

### Ban detection / auto-pause
```
periodicamente por instância (janela móvel, ex: últimos 100 envios):
  taxa_entrega = entregues / enviados
  taxa_erro    = falhas / tentativas
  se taxa_entrega < limite  OU  taxa_erro > limite:
     instance.status = paused; health_score-- ; alerta
```

---

## 6. Drivers por canal (padrão Strategy)

```
interface ChannelDriver:
  send(instance, recipient, content) -> Result
  presence(instance, recipient, duration)     # no-op onde não se aplica
  throttle_defaults() -> ThrottleProfile
  compliance() -> { session_window_hours, requires_optin, needs_unsubscribe }
  capabilities() -> { text, media, subject, buttons }
```

| Canal | Rota de envio | Ritmo / risco | Regra de compliance |
|---|---|---|---|
| **WhatsApp** | evolution-api (texto) + n8n (mídia) | Alto risco · pacing humano | Janela 24h / opt-in ou template. Lista fria = ban |
| **Telegram** | Bot API | ~30 msg/s global, ~1/s por chat | Usuário precisa ter dado `/start` no bot |
| **Instagram** | Meta Graph API via CRM | Muito restrito | Só dentro de 24h ou com *message tag*. Sem blast frio |
| **Facebook** | Meta Graph API via CRM | Muito restrito | Idem: janela 24h + tags. Sem blast frio |

> **Email fora de escopo** (decisão do produto). Fica registrado que, se um dia entrar, o driver
> precisará de SPF/DKIM/DMARC + link de descadastro + warm-up de domínio.

**Realidade de política (importante):** canais Meta (WhatsApp/IG/FB) **não permitem disparo em massa
para lista fria** — é regra da Meta. Servem para **reativar** (quem já falou com você) ou **template
aprovado**. Para *cold/massa* de verdade sem restrição de janela: **Telegram**. A UI deve sinalizar isso.

**Atalho de implementação:** o CRM (base Chatwoot) já entrega em todos esses canais. Os drivers de
Telegram/IG/FB **não reimplementam** o envio — apenas mandam o CRM postar a mensagem de saída na
inbox certa. Só o WhatsApp tem caminho especial (evolution + n8n) por causa do anti-ban.

---

## 7. Seletor de canal (frontend)

Zero backend novo — reusa o que já existe:
- Store: `useAppDataStore()` → `inboxes`, `fetchInboxes()`
- Service: `InboxesService.list()` (endpoint `/inboxes`)
- Tipo `Inbox` já traz: `id, name, channel_type, phone_number, email, avatar_url`

**Multi-seleção (checkbox):** o usuário marca 1 ou várias inboxes → viram o **pool de rotação** da
campanha (gravadas em `cmp_campaign_channels`). O motor lê o `channel_type` de cada inbox → escolhe o
driver e distribui a carga entre os números do pool (rotação + warm-up por número).

> Nota de UX: idealmente o pool é de inboxes do **mesmo `channel_type`** (não misturar WhatsApp com
> Telegram na mesma campanha), pois conteúdo e regras diferem por canal.

---

## 8. Mudança de contrato no n8n (pequena, retrocompatível)

Para enviar de números diferentes (rotação do pool), o payload ganha `instance`:
```json
{ "telefone": "...", "video_url": "...", "media": "...", "instance": "campanha-num-A" }
```
Se `instance` vier vazio → usa o número padrão de hoje (não quebra o fluxo atual).

---

## 9. Compliance
- **Opt-out automático**: palavra-chave ("SAIR"/"PARAR") → entra em `cmp_suppression`, nunca mais recebe.
- **Janela de 24h** (WA/IG/FB): fora dela exige template/tag — sinalizar na UI.
- **Email**: link de descadastro obrigatório + autenticação de domínio (SPF/DKIM/DMARC).
- **LGPD**: rastrear origem/consentimento; suppression respeitada em todas as campanhas.

---

## 10. Observabilidade
Dashboard por campanha e por número: enviados / entregues / lidos / respondidos, health score de cada
número, taxa de opt-out e "termômetro de risco de ban".

---

## 11. Isolamento (o que compartilha com o sistema atual)
- **Novo serviço**, tabelas `cmp_*` novas, **Redis DB /6** dedicado → sem colisão.
- **Leitura** de contatos/inboxes do CRM; **envio** de canais não-WA via API do CRM.
- **WhatsApp**: evolution-api + o **mesmo** webhook n8n de hoje (só ganha campo `instance`).
- Único acoplamento real: o **número do WhatsApp**. Resolvido usando **inbox/número dedicado pra
  campanha** (o usuário escolhe o canal no seletor) → risco de ban do atendimento = zero.

---

## 12. Fases de entrega
1. **MVP**: WhatsApp (via n8n/evolution) + **segmento dinâmico** + **pool de números com rotação** +
   warm-up + jitter + horário comercial + **spintax local** + suppression/opt-out + start manual.
   Sobe desligado; só age ao criar campanha.
2. **Multicanal**: driver Telegram + seletor multi-check ligado (IG/FB dentro da janela depois).
3. **Anti-ban avançado**: afinamento de warm-up/rotação, presença "digitando", A/B de variantes.
4. **Inteligência**: personalização por **LLM** (Python), ban detection + auto-pause + dashboards.

---

## 13. Decisões (definidas)
- **Personalização**: ✅ spintax local no MVP; LLM em fase posterior (Fase 4).
- **Email**: ❌ fora de escopo.
- **Números por campanha**: ✅ pool desde o MVP — seletor multi-check (1..N inboxes) com rotação.
- **Segmentação**: ✅ segmento dinâmico já no MVP (filtro salvo + snapshot no disparo p/ reprodutibilidade).

- **Snapshot**: ✅ no disparo, congela a lista do segmento (reprodutível).

- **Gatilho do MVP**: ✅ start manual (on_label fica pra fase posterior).

---

## 14. Plano de implementação do MVP (por serviço)

Escopo do MVP: **WhatsApp** · segmento dinâmico · pool de números c/ rotação · warm-up + jitter +
horário · spintax local · suppression/opt-out · **start manual**.

### A. evo-campaign-engine (Go) — SERVIÇO NOVO ✅
- [x] Migrations das tabelas `cmp_*`
- [x] Conexão Redis (DB /6): buckets, next_slot, sent_since_pause
- [x] **Scheduler**: promove jobs vencidos (DB poll → channel de trabalho)
- [x] **Worker pool**: consome, aplica throttle, envia, atualiza estado
- [x] **Throttle**: warm-up cap + jitter + quiet hours + token bucket por número
- [x] **Rotação de números** no pool (`cmp_campaign_channels`)
- [x] **Spintax** local (variação de texto)
- [x] **Driver WhatsApp**: texto → evolution-api · mídia → n8n webhook (campo `instance`)
- [x] Checagem de **suppression**
- [x] REST: criar / iniciar / pausar / cancelar campanha + stats
- [x] Dockerfile + entrypoint + CI matrix + docker-stack.yml

### B. evo-crm (Rails) ✅
- [x] Endpoint interno: **resolver segmento dinâmico** (POST /api/v1/campaigns/resolve_segment) → lista de contatos com phone + timezone
- [x] **Opt-out**: capturar palavra-chave (SAIR/PARAR/STOP/etc) em msg recebida → gravar em `cmp_suppression` via Wisper listener
- [x] Autenticação serviço↔serviço (token header X-Campaign-Engine-Token)
- [x] (inboxes já expostas via Rails — reutilizadas)

### C. evo-ai-frontend (React) ✅
- [x] Form de campanha: **seletor multi-check de canal** (reusa `useAppDataStore().inboxes`)
- [x] Editor de conteúdo (texto + mídia URL + spintax suporte)
- [x] Variant management com weight para rotação
- [x] Botão Iniciar/Pausar/Deletar + dashboard básico (enviados/total)

### D. n8n ✅
- [x] Webhook retrocompatível aceita `instance` field (novo) sem quebrar payloads legados (sem instance)

### E. evo-processor (Python)
- [ ] **Nada no MVP** (entra só na Fase 4, com LLM).

### F. Infra (compose / swarm) ✅
- [x] Adicionar serviço `evo-campaign-engine` no `docker-stack.yml` + envs (Redis /6)
- [x] Adicionar na matriz do `docker-build.yml`
