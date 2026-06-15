package coordinator

import (
	"encoding/json"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/nats-io/nats.go"
)

// Coordinator é o roteador principal do Swarm.
type Coordinator struct {
	EventBus *evbus.EventBus
}

// NewCoordinator instancia um novo Roteador.
func NewCoordinator(bus *evbus.EventBus) *Coordinator {
	return &Coordinator{
		EventBus: bus,
	}
}

// Start inicializa o Coordinator e assina os eventos essenciais do NATS.
func (c *Coordinator) Start() error {
	log.Println("🚀 Iniciando Swarm Coordinator...")

	// Assina eventos de novos Leads (Exemplo de Ingress Point para workflows)
	_, err := c.EventBus.Subscribe(string(events.EventLeadCreated), c.handleLeadCreated)
	if err != nil {
		return err
	}

	// Assina eventos de conclusão de agentes (Handoff e Continuidade do Workflow)
	_, err = c.EventBus.Subscribe(string(events.EventAgentFinished), c.handleAgentFinished)
	if err != nil {
		return err
	}

	log.Println("✅ Swarm Coordinator aguardando eventos...")
	return nil
}

// handleLeadCreated é acionado quando um Lead é criado.
// Aqui a Engine de Workflow (DAG) seria instanciada e o primeiro nó (Agente) disparado.
func (c *Coordinator) handleLeadCreated(msg *nats.Msg) {
	log.Printf("[Coordinator] Novo evento recebido: %s", msg.Subject)

	// Na prática:
	// 1. Decodifica o Lead.
	// 2. Busca no banco qual Workflow está atrelado a esse Lead/Sinal.
	// 3. Monta o DAG em memória (workflow.NewWorkflow).
	// 4. Executa o DAG passando o contexto do Lead.
}

// handleAgentFinished intercepta o fim de um trabalho de um agente.
func (c *Coordinator) handleAgentFinished(msg *nats.Msg) {
	log.Printf("[Coordinator] Agente finalizou tarefa: %s", msg.Subject)

	var event events.AgentFinishedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[Coordinator] Erro ao decodificar AgentFinishedEvent: %v", err)
		return
	}

	// Na prática:
	// 1. Busca a Engine de Workflow (DAG) pausada pelo TaskID/TraceID no Redis.
	// 2. Atualiza o estado da DAG com o Result (event.Result).
	// 3. Verifica o DAG para encontrar os próximos nós (NextNodes).
	// 4. Dispara EventAgentStarted para os próximos agentes.
}
