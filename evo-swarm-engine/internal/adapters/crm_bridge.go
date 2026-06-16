package adapters
	
	import (
		"context"
		"encoding/json"
		"fmt"
		"log"
	
		"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/crm"
		evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
		"github.com/nats-io/nats.go"
	)
	
	// CRMBridgeAdapter escuta comandos NATS (Event-Driven) e os converte em chamadas HTTP REST
	// para o CRM Core (Ruby on Rails). Quando o Rails ganhar suporte nativo ao NATS, este adapter pode ser desligado.
	type CRMBridgeAdapter struct {
		EventBus *evbus.EventBus
		Client   *crm.Client
	}
	
	func NewCRMBridgeAdapter(bus *evbus.EventBus) *CRMBridgeAdapter {
		return &CRMBridgeAdapter{
			EventBus: bus,
			Client:   crm.NewClient(),
		}
	}
	
	func (a *CRMBridgeAdapter) Start(ctx context.Context) error {
		log.Println("🌉 [CRMBridge] Iniciado. Traduzindo comandos NATS -> REST API do CRM")
	
		// Escuta comandos de transferência de humano
		_, err := a.EventBus.Conn.QueueSubscribe("crm.command.transfer", "crm_bridge_group", a.handleTransferCommand)
		if err != nil {
			return err
		}
	
		// Futuro: Adicionar crm.command.create_lead, crm.command.update_pipeline, etc.
		return nil
	}
	
	func (a *CRMBridgeAdapter) handleTransferCommand(msg *nats.Msg) {
		var payload struct {
			ConversationID string `json:"conversation_id"`
			Reason         string `json:"reason"`
			TaskID         string `json:"task_id"`
		}
	
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Printf("❌ [CRMBridge] Erro ao decodificar payload de transferência: %v", err)
			msg.Respond([]byte(`{"status":"error", "message":"invalid payload"}`))
			return
		}
	
		log.Printf("🌉 [CRMBridge] Executando REST POST para Transferência. Conv: %s", payload.ConversationID)
	
		// Chamadas REST atômicas
		assignEndpoint := fmt.Sprintf("/conversations/%s/assignments", payload.ConversationID)
		a.Client.Post(context.Background(), assignEndpoint, map[string]interface{}{})
	
		statusEndpoint := fmt.Sprintf("/conversations/%s/toggle_status", payload.ConversationID)
		a.Client.Post(context.Background(), statusEndpoint, map[string]interface{}{
			"status": "open",
		})
	
		log.Printf("✅ [CRMBridge] Conversa %s transferida com sucesso no banco Rails.", payload.ConversationID)
	
	// Responde ao Request-Reply do NATS
	type BridgeResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	responsePayload, err := json.Marshal(BridgeResponse{
		Status:  "success",
		Message: "Transferência realizada com sucesso para a fila de atendimento humano.",
	})
	if err != nil {
		log.Printf("❌ [CRMBridge] Erro ao serializar BridgeResponse: %v", err)
		return
	}
	
	if msg.Reply != "" {
		if err := msg.Respond(responsePayload); err != nil {
			log.Printf("❌ [CRMBridge] Erro ao responder via NATS: %v", err)
		}
	} else {
		if err := a.EventBus.Publish("crm.event.transferred", responsePayload); err != nil {
			log.Printf("❌ [CRMBridge] Erro ao publicar crm.event.transferred: %v", err)
		}
	}
}
