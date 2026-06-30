package workflow
	
	import (
		"context"
		"encoding/json"
		"fmt"
		"log"
		"time"
	
		evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
		"github.com/google/uuid"
		"github.com/nats-io/nats.go"
	)
	
	type DAGState struct {
		WorkflowID string                 `json:"workflow_id"`
		ContextID  string                 `json:"context_id"`
		CurrentNode string                 `json:"current_node"`
		NextNodes  []string               `json:"next_nodes"`
		Variables  map[string]interface{} `json:"variables"`
		Status     string                 `json:"status"` // pending, running, completed, failed
		UpdatedAt  time.Time              `json:"updated_at"`
	}

	// WorkflowPayload é o payload usado para comunicar com os agentes
	type WorkflowPayload struct {
		WorkflowID string                 `json:"workflow_id"`
		ContextID  string                 `json:"context_id"`
		Variables  map[string]interface{} `json:"variables"`
	}
	
	// DAGEngine orquestra a máquina de estados distribuída.
	type DAGEngine struct {
		EventBus *evbus.EventBus
		KV       nats.KeyValue
	}
	
	// NewDAGEngine inicializa o motor de workflows e o Bucket KV no JetStream.
	func NewDAGEngine(bus *evbus.EventBus) (*DAGEngine, error) {
		js, err := bus.Conn.JetStream()
		if err != nil {
			return nil, fmt.Errorf("falha ao conectar JetStream: %v", err)
		}
	
		// Garante que o Bucket KV "WORKFLOWS" existe
		kv, err := js.KeyValue("WORKFLOWS")
		if err != nil {
			// Cria o bucket se não existir
			kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
				Bucket: "WORKFLOWS",
				Description: "Estado dos Workflows em Execução",
				TTL:    24 * time.Hour, // O estado expira em 24h caso não concluído
			})
			if err != nil {
				return nil, fmt.Errorf("falha ao criar KV Store: %v", err)
			}
		}
	
		log.Println("⚙️ [DAG Engine] JetStream KV Bucket 'WORKFLOWS' inicializado.")
	
		return &DAGEngine{
			EventBus: bus,
			KV:       kv,
		}, nil
	}
	
	// StartWorkflow inicia um novo fluxo e salva o estado inicial no NATS KV.
	func (d *DAGEngine) StartWorkflow(ctx context.Context, startNode string, contextID string, vars map[string]interface{}) (string, error) {
		workflowID := uuid.New().String()
	
		state := DAGState{
			WorkflowID: workflowID,
			ContextID:  contextID,
			CurrentNode: startNode,
			Variables:  vars,
			Status:     "running",
			UpdatedAt:  time.Now(),
		}
	
		err := d.SaveState(workflowID, state)
		if err != nil {
			return "", err
		}
	
		// Publica o evento inicial para o respectivo Agente (startNode = agent.crm.task, por exemplo)
		payload, err := json.Marshal(WorkflowPayload{
			WorkflowID: workflowID,
			ContextID:  contextID,
			Variables:  vars,
		})
		if err != nil {
			return "", fmt.Errorf("falha ao serializar payload do workflow: %v", err)
		}
	
		d.EventBus.Publish(startNode, payload)
		log.Printf("🚀 [DAG Engine] Workflow %s iniciado. Nó: %s", workflowID, startNode)
	
		return workflowID, nil
	}
	
	// Transition move o workflow para os próximos nós na DAG.
	func (d *DAGEngine) Transition(workflowID string, output map[string]interface{}, nextNodes []string) error {
		state, err := d.GetState(workflowID)
		if err != nil {
			return err
		}
	
		// Atualiza as variáveis globais do workflow com o resultado do nó anterior
		for k, v := range output {
			state.Variables[k] = v
		}
	
		if len(nextNodes) == 0 {
			state.Status = "completed"
			state.UpdatedAt = time.Now()
			if err := d.SaveState(workflowID, state); err != nil {
				log.Printf("⚠️ [DAG Engine] Erro ao salvar estado de conclusão do workflow %s: %v", workflowID, err)
			}
			log.Printf("✅ [DAG Engine] Workflow %s concluído.", workflowID)
			return nil
		}
	
		state.NextNodes = nextNodes
		state.UpdatedAt = time.Now()
		if err := d.SaveState(workflowID, state); err != nil {
			return fmt.Errorf("falha ao salvar estado de transição do workflow %s: %v", workflowID, err)
		}
	
		// Dispara os próximos nós (Pode disparar vários nós paralelamente se branching for suportado)
		for _, node := range nextNodes {
			payload, err := json.Marshal(WorkflowPayload{
				WorkflowID: workflowID,
				ContextID:  state.ContextID,
				Variables:  state.Variables,
			})
			if err != nil {
				log.Printf("❌ [DAG Engine] Erro ao serializar payload para nó %s: %v", node, err)
				continue
			}
			d.EventBus.Publish(node, payload)
			log.Printf("↪️ [DAG Engine] Workflow %s transição para Nó: %s", workflowID, node)
		}
	
		return nil
	}
	
	// SaveState persiste o estado no JetStream KV
	func (d *DAGEngine) SaveState(workflowID string, state DAGState) error {
		data, err := json.Marshal(state)
		if err != nil {
			return err
		}
		_, err = d.KV.Put(workflowID, data)
		return err
	}
	
	// GetState recupera o estado do JetStream KV
	func (d *DAGEngine) GetState(workflowID string) (DAGState, error) {
		var state DAGState
		entry, err := d.KV.Get(workflowID)
		if err != nil {
			return state, err
		}
		
		err = json.Unmarshal(entry.Value(), &state)
		return state, err
	}
