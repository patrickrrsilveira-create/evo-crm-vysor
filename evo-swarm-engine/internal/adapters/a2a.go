package adapters

import (
	"encoding/json"
	"log"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type A2AAdapter struct {
	EventBus *evbus.EventBus
}

func NewA2AAdapter(bus *evbus.EventBus) *A2AAdapter {
	return &A2AAdapter{EventBus: bus}
}

func (a *A2AAdapter) RegisterRoutes(app *fiber.App) {
	a2aGroup := app.Group("/api/v1/a2a")

	// Rota Genérica de A2A (JSON-RPC)
	a2aGroup.Post("/:agent_id", func(c *fiber.Ctx) error {
		agentID := c.Params("agent_id")

		var req models.A2ARequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(models.A2AResponse{
				JSONRPC: "2.0",
				Error:   &models.A2AError{Code: -32700, Message: "Parse error"},
			})
		}

		log.Printf("🤖 [A2A Protocol] Recebida chamada para %s (Method: %s)", agentID, req.Method)

		switch req.Method {
		case "message/send":
			return a.handleMessageSend(c, agentID, req)
		case "message/stream":
			// TODO: Implementar SSE Starlette style para stream real
			return a.handleMessageSend(c, agentID, req) // Fallback para send por enquanto
		default:
			return c.Status(400).JSON(models.A2AResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &models.A2AError{Code: -32601, Message: "Method not found"},
			})
		}
	})
}

func (a *A2AAdapter) handleMessageSend(c *fiber.Ctx, agentID string, req models.A2ARequest) error {
	// Extrai texto da request A2A
	var content string
	if parts, ok := req.Params["parts"].([]interface{}); ok && len(parts) > 0 {
		if part, ok := parts[0].(map[string]interface{}); ok {
			if text, ok := part["text"].(string); ok {
				content = text
			}
		}
	}

	// 1. Dispara o evento de Input no Barramento para a Engine processar
	// (Simula o que o webhook faria, ativando a DAG)
	eventData, _ := json.Marshal(map[string]interface{}{
		"source":   "a2a_protocol",
		"content":  content,
		"agent_id": agentID,
	})
	a.EventBus.Publish(string(events.EventMessageReceived), eventData)

	// Na implementação síncrona/SSE real, aqui esperaríamos um channel de resposta.
	// Por ser um MVP de paridade, retornamos o TaskID aceito (Pending state)

	taskID := uuid.New().String()

	task := models.A2ATask{
		ID:        taskID,
		ContextID: req.Params["contextId"].(string),
		Status: models.A2ATaskStatus{
			State:     "pending",
			Timestamp: time.Now(),
		},
		Kind: "task",
	}

	return c.JSON(models.A2AResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  task,
	})
}
