package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/memory"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// MemoryService abstrai o acesso ao Redis e PostgreSQL via NATS
type MemoryService struct {
	EventBus    *evbus.EventBus
	ShortMemory *memory.ShortMemoryManager
	db          *gorm.DB
}

func NewMemoryService(bus *evbus.EventBus, db *gorm.DB, redisClient *redis.Client) *MemoryService {
	return &MemoryService{
		EventBus:    bus,
		ShortMemory: memory.NewShortMemoryManager(redisClient),
		db:          db,
	}
}

func (s *MemoryService) Start() error {
	log.Println("💾 [MemoryService] Iniciado. Três Camadas de Memória Ativas.")

	// Serviço de leitura (Request-Reply)
	_, err := s.EventBus.Conn.QueueSubscribe("service.memory.query", "memory_service_group", s.handleQuery)
	if err != nil {
		return err
	}

	// Serviço de escrita passiva para IA (Outbound)
	_, err = s.EventBus.Conn.QueueSubscribe("outbound.message", "memory_logger_group", s.handleOutbound)
	if err != nil {
		return err
	}

	// Serviço de escrita passiva para Usuários (Inbound)
	_, err = s.EventBus.Conn.QueueSubscribe(string(events.EventMessageReceived), "memory_logger_group", s.handleInbound)
	return err
}

func (s *MemoryService) handleQuery(msg *nats.Msg) {
	var req struct {
		ConversationID string `json:"conversation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		log.Printf("❌ [MemoryService] Erro ao decodificar request de query: %v", err)
		msg.Respond([]byte(`{"error": "invalid request format"}`))
		return
	}

	// Consulta apenas a Short Term Memory (Redis) para a LLM
	history, err := s.ShortMemory.GetRecentHistory(context.Background(), req.ConversationID)
	if err != nil {
		log.Printf("❌ [MemoryService] Erro ao consultar histórico Redis: %v", err)
		msg.Respond([]byte(`{"error": "failed to get memory"}`))
		return
	}

	resp, err := json.Marshal(struct {
		History []models.LLMMessage `json:"history"`
	}{
		History: history,
	})
	if err != nil {
		log.Printf("❌ [MemoryService] Erro ao serializar history: %v", err)
		msg.Respond([]byte(`{"error": "failed to serialize memory"}`))
		return
	}
	
	msg.Respond(resp)
}

// Payload padrão para mensagens trocadas
type MessagePayload struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func (s *MemoryService) handleOutbound(msg *nats.Msg) {
	var payload MessagePayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Printf("❌ [MemoryService] Erro ao decodificar outbound message: %v", err)
		return
	}

	sender := payload.Sender // Na versão atual, sender armazena o conversationID
	content := payload.Content

	if sender != "" && content != "" {
		// 1. Salva no Redis (Short Term Memory)
		if err := s.ShortMemory.AddMessageToHistory(context.Background(), sender, models.LLMMessage{
			Role:    "assistant",
			Content: content,
		}); err != nil {
			log.Printf("❌ [MemoryService] Erro ao salvar outbound no Redis: %v", err)
		}

		// 2. Salva no PostgreSQL (Conversation History / Long Term)
		if s.db != nil {
			if err := s.db.Create(&models.ConversationMessage{
				ConversationID: sender,
				Role:           "assistant",
				Content:        content,
			}).Error; err != nil {
				log.Printf("❌ [MemoryService] Erro ao salvar outbound no PostgreSQL: %v", err)
			}
		}
	}
}

func (s *MemoryService) handleInbound(msg *nats.Msg) {
	var payload MessagePayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Printf("❌ [MemoryService] Erro ao decodificar inbound message: %v", err)
		return
	}

	sender := payload.Sender
	content := payload.Content

	if sender != "" && content != "" {
		// 1. Salva no Redis (Short Term Memory)
		if err := s.ShortMemory.AddMessageToHistory(context.Background(), sender, models.LLMMessage{
			Role:    "user",
			Content: content,
		}); err != nil {
			log.Printf("❌ [MemoryService] Erro ao salvar inbound no Redis: %v", err)
		}

		// 2. Salva no PostgreSQL (Conversation History / Long Term)
		if s.db != nil {
			if err := s.db.Create(&models.ConversationMessage{
				ConversationID: sender,
				Role:           "user",
				Content:        content,
			}).Error; err != nil {
				log.Printf("❌ [MemoryService] Erro ao salvar inbound no PostgreSQL: %v", err)
			}
		}
	}
}
