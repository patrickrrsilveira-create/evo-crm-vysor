# Master Plan: Re-Arquitetura Evo AI Processor (Project Go-Swarm)

Abaixo está o plano estrutural definitivo e prioritário para recriar o coração da plataforma, abandonando as "gambiarras" atuais, escrevendo do zero em **Go (Golang)**, implementando **Agent Swarms (Enxame de Agentes Real)** e dominando a orquestração Omni-channel.

## Princípios da Nova Arquitetura
1. **Zero Gambiarras**: Código tipado (Go), alta concorrência nativa (Goroutines) e arquitetura orientada a eventos.
2. **Aproveitamento Total das Telas**: O Frontend (Vue/Nuxt/React) atual não sofre 1 linha de alteração. O novo backend em Go "imita" 100% os contratos da API original.
3. **Orquestração Independente**: A IA se torna o "Cérebro do Meio", comunicando diretamente com Evolution API, WA Cloud, Z-API, etc. O Chatwoot vira apenas um "visor" para humanos.

---

## 1. Mapeamento e Migração da Camada de Dados (Database)
**Prioridade:** `Altíssima` (O alicerce do sistema)
A primeira etapa é portar a base PostgreSQL do Python (`SQLAlchemy`) para o Go.
* **Ferramenta:** Usaremos o **GORM** ou **Ent** (Facebook) para modelar o banco no Go.
* **Tabelas Migradas 1 pra 1:** 
  - `evo_core_agents`, `evo_core_agent_folders`, `evo_core_api_keys`
  - `evo_core_custom_tools`, `evo_core_custom_mcp_servers`
  - `evo_ai_agent_processor_sessions`, `evo_core_agent_integrations`
* **Vantagem:** Preservamos os dados de todos os clientes atuais sem perda ou necessidade de migração estrutural nas telas.

## 2. Reescrita dos Endpoints da API (Frontend Compatibility)
**Prioridade:** `Alta`
O painel visual do Evo precisa continuar funcionando.
* **Ferramenta:** Usaremos o framework web **Fiber** (Go), que é extremamente rápido e parecido com o Express/FastAPI.
* **Rotas a serem recriadas:**
  - `GET/POST /api/v1/agents/*` (Para as telas de criação de agentes)
  - `GET/POST /api/v1/integrations/*` (Para as telas de integrações)
  - `GET/POST /api/v1/tools/*` (Para MCP e ferramentas customizadas)
* **Resultado:** O frontend fará o *request* para a mesma URL, mas o binário em Go processará 100x mais rápido que o Python.

## 3. Construção do Engine de Agent Swarm (O Fim do Proxy)
**Prioridade:** `Crítica` (O core business)
Abandonaremos o `google.adk` (que força o modelo Proxy/Tool Calling) e construiremos nosso próprio *Swarm Router* em Go (baseado no conceito do *LangChainGo* ou *Actor Model*).
* **Como vai funcionar:** 
  1. O **Router Agent** (Gabriela) recebe a mensagem do usuário.
  2. O Router decide que a especialidade é do Roberto.
  3. Ela faz um **Handoff (Transferência de Contexto)** para o Roberto (Passa a variável de memória do usuário para a Goroutine do Roberto).
  4. O Roberto processa a API, gera o áudio, e **ele mesmo** envia a resposta para o WhatsApp do cliente.
* **Vantagem:** Fim das "caixas pretas" de sub-agentes escondidos. Cada agente opera de forma autônoma e assíncrona.

## 4. O Sistema "Omni-channel Core"
**Prioridade:** `Alta`
A IA não dependerá mais de receber webhooks apenas do Chatwoot.
* O motor em Go terá "Adaptadores" nativos (Interfaces em Go):
  - `type ChannelAdapter interface { SendMessage(), ReadWebhook() }`
* **Integrações Nativas:**
  - `EvolutionAPIAdapter` (Conexão nativa com a api Evolution Go/Node)
  - `WhatsAppCloudAdapter` (Oficial Meta)
  - `WApiAdapter` / `ZApiAdapter`
* **A Magia:** A IA recebe a mensagem direto do WhatsApp, responde em 1 segundo e manda um "espelho" pro Chatwoot só pra ficar registrado caso um humano queira ler depois.

## 5. Solucionando a "Parede do Chatwoot" (Perfeição dos Avatares)
**Prioridade:** `Média`
Já que o Chatwoot exige um usuário/bot no banco de dados para renderizar a foto na tela:
* **A Solução Definitiva (Sem Forks):**
  - Quando você cria um Agente no painel do Evo (ex: "Roberto"), o motor em Go faz uma chamada secreta na API de Administração do Chatwoot (`POST /api/v1/accounts/1/agent_bots`).
  - O Go cria automaticamente o Bot do Roberto lá dentro do Chatwoot com a foto dele.
  - Quando ocorre o "Handoff" no Swarm, o motor Go posta a resposta no Chatwoot usando o token daquele Bot específico!
* **Resultado:** Na tela do Chatwoot, você verá a foto e o nome de cada Agente separadamente na mesma conversa. Perfeição absoluta, sem gambiarras de texto e sem precisar modificar o código do Chatwoot.

## 6. A Migração das Integrações (TTS, Webhooks e Provedores de IA)
Respondendo à sua dúvida sobre aproveitar as integrações atuais:
* **As Configurações e o Banco de Dados:** Serão **100% aproveitados**. As chaves de API do ElevenLabs, OpenAI, webhooks salvos, e configurações no painel continuarão no mesmo lugar.
* **O Código das Integrações:** Como Go e Python são linguagens diferentes, o código bruto (as linhas de código) não é copiado, ele é **re-implementado em Go**.
  - **A Boa Notícia:** Fazer integrações de API (Webhooks, TTS, ChatGPT) em Go é extremamente simples e nativo. O Go lida com I/O de rede (HTTP) de forma muito mais rápida que o Python.
  - **Microserviço (Plano B):** Se houver alguma biblioteca do Python muito exclusiva que você não queira reescrever, podemos manter o Python vivo apenas como um "Worker de Ferramentas". O Go orquestra tudo na velocidade da luz, e quando precisa rodar uma ferramenta muito complexa, ele joga pro Python processar e devolver. Mas, para TTS, Webhooks e LLMs padrão, reescrever em Go é o caminho mais limpo, seguro e performático.

## 7. A Infraestrutura Pesada (Redis, MinIO, RabbitMQ)
Você pontuou muito bem. O Evo não é só banco de dados, tem uma infraestrutura pesada rodando em volta. No Go, nós lidaremos com isso de forma ainda mais performática:
* **Redis:** O Go possui a biblioteca `go-redis`, que é o padrão da indústria. Continuaremos usando o seu Redis atual para controle de filas, cache de sessões do WhatsApp e controle de Handoff entre os agentes.
* **MinIO (S3):** Todos os áudios gerados pelo TTS, imagens e documentos que hoje o Python salva no MinIO, o Go salvará usando o SDK oficial `minio-go`. O seu *bucket* atual de arquivos continuará intacto. Os agentes do Swarm farão upload/download em milissegundos usando streams do Go.

## 8. Integrações Avançadas (OAuth, MCP Servers e Apps)
As telas que você mostrou de OAuth (Google, Microsoft), Webhooks e MCP Servers funcionam sob protocolos universais.
* **OAuth 2.0:** O processo de login (Callback URL) para Google Calendar, Sheets e Teams será reescrito nas rotas do Fiber (Go). O painel Vue continua enviando os dados pro backend, e o Go gerencia os *Tokens* no banco.
* **Servidores MCP (Model Context Protocol):** O Go possui bibliotecas nativas para gerenciar processos MCP (via `stdio` ou `SSE`). O Go vai "ligar" o seu servidor MCP do Notion, Github e Stripe, mantendo as ferramentas ativas para os Agentes do Enxame consumirem.

## 9. Estratégia de Transição Suave (Blue/Green Deployment)
Como testar tudo isso sem parar a operação atual (Zero Downtime)?
1. **O Motor Gêmeo:** Criaremos o motor Go rodando em paralelo no seu servidor, mas em uma porta diferente (ex: Python na `8000` e Go na `8001`).
2. **Mesma Base, Duas Cabeças:** O motor Go se conectará no **mesmo** banco de dados e no mesmo Redis da produção. Ele consegue ler e escrever os mesmos dados sem corromper nada.
3. **Ambiente de Homologação (Staging):** Subiremos uma cópia do seu painel Vue apontando secretamente para a porta `8001` (Go). A sua equipe poderá criar agentes, testar integrações e bater webhooks na API do Go, validando que o Swarm está perfeito. Os clientes reais continuam sendo atendidos pelo Python.
4. **A Virada de Chave:** Quando o Go estiver 100% testado e validado, nós apenas mudamos o roteamento do servidor (Traefik/Nginx): mandamos a porta principal apontar para o Go e desligamos o Python. O cliente final não percebe sequer 1 segundo de queda.

## User Review Required

> [!IMPORTANT]
> Patrick, este é um projeto de re-arquitetura profundo (papo de nível Sênior/Staff Engineer). É um produto completamente novo "por baixo do capô", mas que mantém toda a casca de ouro que você já tem no frontend e no banco de dados, além de reaproveitar toda a infraestrutura (Redis, MinIO).
> 
> Esse é o plano nível 1. Faz sentido iniciarmos a execução desse escopo em fases? A Fase 1 seria iniciar os modelos do banco no novo projeto em Go. Como deseja prosseguir?
