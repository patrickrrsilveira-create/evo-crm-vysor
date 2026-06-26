# Plano de Execução: Motor Go Swarm (Python -> Go)

## Visão Estratégica
A abordagem pragmática: **Vamos trocar apenas o motor do carro**. 
O Frontend (Vue), o CRM/Inbox (Chatwoot em Ruby) e a Base de Dados (PostgreSQL) permanecem **intactos**. O alvo da cirurgia é arrancar o `evo-ai-processor-community` (Python) e substituí-lo por um **Swarm Engine em Go**, mas construído com padrões de software Enterprise (Arquitetura orientada a eventos e microsserviços internos).

---

## Fase 1: Paridade de API e Banco de Dados (O Esqueleto)
O objetivo desta fase é fazer com que o painel do Evo (Frontend) consiga salvar e ler configurações no novo motor em Go, sem perceber que o Python sumiu.

* **Conexão com PostgreSQL:** Usar **GORM** no Go para conectar no seu banco atual e mapear as tabelas exatas (`evo_core_agents`, `evo_core_custom_tools`, etc.).
* **Recriar API REST (Fiber):** Criar as rotas que o painel Vue consome hoje (ex: `/api/v1/agents`, `/api/v1/tools`) usando **Fiber** (Go structs).
* **Resultado Esperado:** O painel do Evo funciona 100% lendo e gravando no banco usando o Go.

## Fase 2: O Motor Interno Orientado a Eventos (A Fundação Go)
Nesta etapa substituímos as bibliotecas engessadas do Python (FastAPI, Pydantic, CrewAI) por tecnologias nativas e de alta performance no ecossistema Go.

* **Event Bus (NATS):** Toda a comunicação interna será orientada a eventos usando **NATS**. Ao invés de funções chamando umas as outras de forma blocante, teremos eventos como `conversation.created`, `agent.started`, `agent.finished`.
* **Memory Engine Dedicado:** Ao invés de espalhar estado em variáveis, criaremos um serviço interno robusto (`memory-service`):
  - **Curto Prazo:** Redis (Contexto da conversa atual).
  - **Médio Prazo:** PostgreSQL (Auditoria e Histórico).
  - **Longo Prazo/Semântica:** pgvector (RAG e Base de Conhecimento).

## Fase 3: Swarm Coordinator & Workflow Engine (A Inteligência)
Aqui matamos o LangGraph e o Google ADK para criar o verdadeiro Enxame Nativo.

* **Workflow Engine Nativo (DAG):** Um motor de grafos direcional (DAG) construído em Go puro (Node -> Action -> Decision -> Parallel -> Join). O fluxo da conversa não fica preso a frameworks de terceiros.
* **Swarm Coordinator:** Um roteador central que não executa ferramentas, apenas **recebe a tarefa, planeja e distribui** o trabalho.
* **Specialist Agents:** Agentes independentes rodando em *Goroutines* que assinam (subscribe) tópicos no NATS. Teremos: *CRM Agent, Voice Agent, WhatsApp Agent, Knowledge Agent*. O Coordinator joga o evento no NATS, o especialista certo pega e resolve.

## Fase 4: Integrações e MCP Service (Os Músculos)
* **Serviço MCP Dedicado:** Um módulo interno (`mcp-service`) inteiramente responsável por *Tool Discovery*, *Tool Execution* e *Tool Registry*, mantendo os servidores externos (Notion, Github, Stripe) isolados da lógica de conversação.
* **Ferramentas Nativas:** Recriar Webhooks, Calendar e integrações cruciais de forma nativa no Go.

## Fase 5: Perfeição do Avatar e Omni-Channel (O "Pulo do Gato")
* **Go -> Evolution API:** O Go vai ler a mensagem e enviar a resposta **diretamente para a Evolution API** (bypassando a limitação do Chatwoot). Isso garante que o cliente no WhatsApp receba a mensagem com a foto e o nome exato do Agente Especialista.
* **Mirroring no Chatwoot:** Em milissegundos, o Go escreve uma cópia invisível dessa mensagem direto no banco do Chatwoot, apenas para que a sua equipe humana consiga ler o histórico na tela.

## Fase 6: Transição Zero-Downtime (Blue/Green)
* **Passo A:** O Motor Go sobe na porta `8001` (ao lado do Python na `8000`).
* **Passo B:** Homologação interna testando o Go (criando agentes, gerando áudio, conversando).
* **Passo C:** Alteramos o roteador (Nginx/Traefik) para mandar o tráfego oficial para a porta `8001`.
* **Passo D:** Desligamos o container Python para sempre.

---

## User Review Required

> [!IMPORTANT]
> Patrick, este é o plano final consolidado. Ele preserva a abordagem inteligente de manter o Chatwoot e o Frontend intactos (salvando meses de trabalho), mas garante que a **estrutura interna do novo motor em Go seja de nível Enterprise** (NATS, Workflow DAG, MCP Nativo).
>
> Ao aprovar, nossa primeira missão (Fase 1) será inicializar o repositório em Go (`evo-swarm-engine`) e configurar a conexão (GORM) com as tabelas do seu banco de dados atual. Podemos prosseguir com esse planejamento?
