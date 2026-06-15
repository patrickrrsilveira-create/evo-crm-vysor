# Plano de Execução: Motor Go Swarm (Python -> Go)

## Visão Estratégica
A abordagem pragmática: **Vamos trocar apenas o motor do carro**. 
O Frontend (Vue), o CRM/Inbox (Chatwoot em Ruby) e a Base de Dados (PostgreSQL) permanecem **intactos**. O alvo da cirurgia é arrancar o `evo-ai-processor-community` (Python) e substituí-lo por um **Swarm Engine em Go**.

---

## Fase 1: Paridade de API e Banco de Dados (O Esqueleto)
O objetivo desta fase é fazer com que o painel do Evo (Frontend) consiga salvar e ler configurações no novo motor em Go, sem perceber que o Python sumiu.

* **Conexão com PostgreSQL:** Usar **GORM** no Go para conectar no seu banco atual e mapear as tabelas exatas (`evo_core_agents`, `evo_core_custom_tools`, etc.).
* **Recriar API REST (Fiber):** Criar as rotas que o painel Vue consome hoje:
  - `GET/POST /api/v1/agents`
  - `GET/POST /api/v1/integrations`
  - `GET/POST /api/v1/tools`
* **Resultado Esperado:** O painel do Evo funciona 100% lendo e gravando no banco usando o Go.

## Fase 2: O Motor Swarm Nativo (O Cérebro)
Aqui abandonamos o `google.adk` e o `LangGraph` (que forçam a arquitetura de "agente escondido/proxy") e criamos a inteligência real de Enxame (Swarm) usando Goroutines de altíssima performance.

* **Coordinator Router:** A IA principal (ex: Gabriela) recebe a mensagem, analisa a intenção e faz o *Handoff* (transferência).
* **Handoff em Memória:** Ao invés de fazer chamadas de rede lentas, o Go transfere a memória da conversa diretamente para a *Goroutine* do Agente Especialista (ex: Roberto).
* **Processamento Paralelo:** O Agente Especialista executa as ferramentas (TTS, Busca) de forma isolada e ultra-rápida.

## Fase 3: Integrações e Infraestrutura (Os Músculos)
O Go precisa reconectar com tudo o que o Python falava hoje.

* **MinIO e Redis:** Conectar o Go no MinIO (para salvar os áudios gerados pelo TTS) e no Redis (para cache de sessões).
* **Ferramentas Nativas:** Recriar as funções de Webhook, envio de e-mail e geração de áudio (TTS ElevenLabs/Fish/OpenRouter) em Go.
* **Model Context Protocol (MCP):** Implementar cliente MCP em Go para que a IA consiga continuar conversando com seus servidores externos (Notion, Github, Stripe).

## Fase 4: Perfeição do Avatar e Omni-Channel (O "Pulo do Gato")
Aqui nós resolvemos o problema da foto e do nome do robô (que o Chatwoot bloqueava).

* **Go -> Evolution API:** O Go vai ler a mensagem e enviar a resposta **diretamente para a Evolution API** (bypassando a limitação do Chatwoot). Isso garante que o cliente no WhatsApp receba a mensagem com a foto e o nome exato do Agente Especialista (Roberto ou Gabriela).
* **Mirroring no Chatwoot:** Em milissegundos, o Go escreve uma cópia invisível dessa mensagem direto no banco do Chatwoot, apenas para que a sua equipe humana consiga ler o histórico na tela.
* **Status "Digitando/Gravando":** O Go emitirá os eventos de *Presence* perfeitamente controlados direto na API da Evolution.

## Fase 5: Transição Zero-Downtime (Blue/Green)
* **Passo A:** O Motor Go sobe na porta `8001` (ao lado do Python na `8000`).
* **Passo B:** Homologação interna testando o Go (criando agentes, gerando áudio, conversando).
* **Passo C:** Alteramos o roteador (Nginx/Traefik) para mandar o tráfego oficial para a porta `8001`.
* **Passo D:** Desligamos o container Python para sempre.

---

## User Review Required

> [!IMPORTANT]
> Patrick, este é um plano cirúrgico, estruturado e realista para trocar a inteligência do sistema sem reinventar a roda do painel.
>
> Ao aprovar, nossa primeira missão (Fase 1) será inicializar o repositório em Go (`evo-swarm-engine`) e configurar a conexão (GORM) com as tabelas do seu banco de dados atual. Podemos prosseguir com esse planejamento?
