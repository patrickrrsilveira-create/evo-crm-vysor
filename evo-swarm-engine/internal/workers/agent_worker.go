package workers

import (
	"encoding/json"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/nats-io/nats.go"
)

// AgentWorker é o handler isolado que escuta e processa tarefas para um agente
type AgentWorker struct {
	EventBus *evbus.EventBus
}

// NewAgentWorker instancia um novo Worker.
func NewAgentWorker(bus *evbus.EventBus) *AgentWorker {
	return &AgentWorker{
		EventBus: bus,
	}
}

// Start inicia a escuta de eventos direcionados a Agentes.
func (w *AgentWorker) Start() error {
	log.Println("🤖 Agent Worker iniciado, aguardando tarefas...")

	// Inscreve-se usando QueueGroup do NATS para balanceamento de carga real!
	// Se tivermos 10 instâncias do motor Go, o NATS vai entregar a tarefa de forma balanceada.
	_, err := w.EventBus.Conn.QueueSubscribe(string(events.EventAgentStarted), "agent_workers_group", w.handleAgentTask)
	return err
}

// handleAgentTask processa a tarefa recebida.
func (w *AgentWorker) handleAgentTask(msg *nats.Msg) {
	var event events.AgentStartedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[AgentWorker] Erro ao decodificar AgentStartedEvent: %v", err)
		return
	}

	log.Printf("[AgentWorker] Recebi Tarefa: [%s] Agente: %s", event.TaskID, event.AgentName)

	// Aqui ocorre a execução (chamada à LLM, LangChain Go, MCP, etc)
	// Após a execução (simulada ou real), nós devolvemos a resposta ao barramento

	finishedEvent := events.AgentFinishedEvent{
		BaseEvent: events.BaseEvent{
			EventID:   event.EventID, // Propaga o ID para rastro
			EventType: events.EventAgentFinished,
			TraceID:   event.TraceID,
		},
		AgentID: event.AgentID,
		TaskID:  event.TaskID,
		Result:  "Processado com Sucesso",
		Success: true,
	}

	data, _ := json.Marshal(finishedEvent)
	w.EventBus.Publish(string(events.EventAgentFinished), data)

	log.Printf("[AgentWorker] Tarefa [%s] concluída e notificada de volta ao NATS", event.TaskID)
}
