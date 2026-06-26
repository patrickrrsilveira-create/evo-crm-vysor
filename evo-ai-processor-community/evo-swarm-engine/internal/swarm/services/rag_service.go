package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/knowledge"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// RAGService abstrai o PgVector via NATS Request-Reply
type RAGService struct {
	EventBus   *evbus.EventBus
	RAGService *knowledge.RAGService
}

func NewRAGService(bus *evbus.EventBus, db *gorm.DB) *RAGService {
	return &RAGService{
		EventBus:   bus,
		RAGService: knowledge.NewRAGService(db),
	}
}

func (s *RAGService) Start() error {
	log.Println("📚 [RAGService] Iniciado. Aguardando consultas NATS (service.rag.query)...")

	_, err := s.EventBus.Conn.QueueSubscribe("service.rag.query", "rag_service_group", s.handleQuery)
	return err
}

func (s *RAGService) handleQuery(msg *nats.Msg) {
	var req struct {
		AgentID string `json:"agent_id"`
		APIKey  string `json:"api_key"`
		Query   string `json:"query"`
		TopK    int    `json:"top_k"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		log.Printf("❌ [RAGService] Erro ao decodificar request de query: %v", err)
		msg.Respond([]byte(`{"error": "invalid request format"}`))
		return
	}

	var parsedAgentID uuid.UUID
	if req.AgentID != "" {
		if parsed, err := uuid.Parse(req.AgentID); err == nil {
			parsedAgentID = parsed
		}
	}

	// Passa a chave de API dinâmica
	apiKeyObj := models.APIKey{Key: req.APIKey}

	contextData, err := s.RAGService.RetrieveContext(context.Background(), parsedAgentID, apiKeyObj, req.Query)
	if err != nil {
		log.Printf("❌ [RAGService] Erro ao consultar contexto RAG: %v", err)
		msg.Respond([]byte(`{"error": "failed to retrieve context"}`))
		return
	}

	resp, err := json.Marshal(struct {
		Context string `json:"context"`
	}{
		Context: contextData,
	})
	if err != nil {
		log.Printf("❌ [RAGService] Erro ao serializar contexto: %v", err)
		msg.Respond([]byte(`{"error": "failed to serialize context"}`))
		return
	}
	
	msg.Respond(resp)
}
