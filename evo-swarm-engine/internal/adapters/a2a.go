package adapters

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type A2AAdapter struct {
	EventBus *evbus.EventBus
}

func NewA2AAdapter(bus *evbus.EventBus) *A2AAdapter {
	return &A2AAdapter{EventBus: bus}
}

func (a *A2AAdapter) RegisterRoutes(app *fiber.App, db *gorm.DB) {
	a2aGroup := app.Group("/api/v1/a2a")

	// Rota Genérica de A2A (JSON-RPC ou Chatwoot Webhook)
	a2aGroup.Post("/:agent_id", func(c *fiber.Ctx) error {
		agentID := c.Params("agent_id")

		// Lê o raw body para sabermos se é Chatwoot ou JSON-RPC
		bodyBytes := c.Body()

		var raw map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &raw); err == nil {
			// Verifica se é um Webhook do Chatwoot (AgentBot)
			if eventType, ok := raw["event"].(string); ok && eventType == "message_created" {
				// É do Chatwoot!
				log.Printf("🤖 [A2A Protocol] Interceptado Webhook do Chatwoot no endpoint A2A para agent_id=%s", agentID)
				
				isIncoming := false
				switch mt := raw["message_type"].(type) {
				case float64:
					isIncoming = mt == 0
				case string:
					isIncoming = strings.EqualFold(mt, "incoming")
				}

				if isIncoming {
					content := ""
					if c, ok := raw["content"].(string); ok {
						content = c
					}

					var conversationID int64
					var accountID int64
					if conv, ok := raw["conversation"].(map[string]interface{}); ok {
						if displayID, ok := conv["display_id"].(float64); ok {
							conversationID = int64(displayID)
						} else if id, ok := conv["id"].(float64); ok {
							conversationID = int64(id)
						}
						if accID, ok := conv["account_id"].(float64); ok {
							accountID = int64(accID)
						}
					}

					sender := ""
					if senderObj, ok := raw["sender"].(map[string]interface{}); ok {
						if name, ok := senderObj["name"].(string); ok {
							sender = name
						}
						if id, ok := senderObj["id"].(float64); ok && sender == "" {
							sender = fmt.Sprintf("user_%d", int64(id))
						}
					}

					eventData, _ := json.Marshal(map[string]interface{}{
						"source":          "chatwoot",
						"content":         content,
						"agent_id":        agentID,
						"sender":          sender,
						"conversation_id": conversationID,
						"account_id":      accountID,
						"payload":         raw,
					})

					a.EventBus.Publish(string(events.EventMessageReceived), eventData)
				}
				return c.SendStatus(200)
			}
		}

		// Se não for Chatwoot, aplica autenticação A2A padrão e processa JSON-RPC
		authHandler := middleware.EvoAuthMiddleware(db)
		err := authHandler(c)
		if err != nil {
			return err
		}

		var req models.A2ARequest
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
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
			return a.handleMessageStream(c, agentID, req)
		default:
			return c.Status(400).JSON(models.A2AResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &models.A2AError{Code: -32601, Message: "Method not found"},
			})
		}
	})
}

// handleMessageSend lida com a versão síncrona. No entanto, enviaremos via NATS para o worker
// processar de forma não bloqueante e retornaremos o "Pending Task", respeitando a arquitetura assíncrona.
func (a *A2AAdapter) handleMessageSend(c *fiber.Ctx, agentID string, req models.A2ARequest) error {
	var content string
	var contextID string

	if ctxID, ok := req.Params["contextId"].(string); ok {
		contextID = ctxID
	} else {
		contextID = uuid.New().String()
	}

	if parts, ok := req.Params["parts"].([]interface{}); ok && len(parts) > 0 {
		if part, ok := parts[0].(map[string]interface{}); ok {
			if text, ok := part["text"].(string); ok {
				content = text
			}
		}
	}

	taskID := uuid.New().String()

	// 1. Dispara o evento de Input no Barramento para a Engine processar
	eventData, _ := json.Marshal(map[string]interface{}{
		"source":     "a2a_protocol",
		"content":    content,
		"agent_id":   agentID,
		"task_id":    taskID,
		"context_id": contextID,
		"is_stream":  false,
	})
	a.EventBus.Publish(string(events.EventMessageReceived), eventData)

	task := models.A2ATask{
		ID:        taskID,
		ContextID: contextID,
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

// handleMessageStream utiliza Server-Sent Events (SSE) nativo para prover respostas em stream hiper-rápidas
// substituindo o bloqueio do Python.
func (a *A2AAdapter) handleMessageStream(c *fiber.Ctx, agentID string, req models.A2ARequest) error {
	var content string
	var contextID string

	if ctxID, ok := req.Params["contextId"].(string); ok {
		contextID = ctxID
	} else {
		contextID = uuid.New().String()
	}

	if parts, ok := req.Params["parts"].([]interface{}); ok && len(parts) > 0 {
		if part, ok := parts[0].(map[string]interface{}); ok {
			if text, ok := part["text"].(string); ok {
				content = text
			}
		}
	}

	taskID := uuid.New().String()

	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	// Dispara o evento indicando que este é um fluxo de stream
	eventData, _ := json.Marshal(map[string]interface{}{
		"source":     "a2a_protocol",
		"content":    content,
		"agent_id":   agentID,
		"task_id":    taskID,
		"context_id": contextID,
		"is_stream":  true,
	})
	a.EventBus.Publish(string(events.EventMessageReceived), eventData)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Inscreve num canal de NATS exclusivo para esta task (Subject: stream.taskID)
		subject := fmt.Sprintf("stream.%s", taskID)
		streamChan := make(chan *nats.Msg, 100)

		sub, err := a.EventBus.Conn.ChanSubscribe(subject, streamChan)
		if err != nil {
			log.Printf("❌ Erro ao assinar canal NATS SSE: %v", err)
			return
		}
		defer sub.Unsubscribe()

		log.Printf("📡 SSE Stream aberto para task %s", taskID)

		for {
			select {
			case msg := <-streamChan:
				var payload map[string]interface{}
				json.Unmarshal(msg.Data, &payload)

				// Se for o marcador de fim, fecha o stream
				if status, ok := payload["status"].(string); ok && status == "completed" {
					log.Printf("📡 SSE Stream %s finalizado.", taskID)
					return
				}

				// Envia o chunk para o cliente
				chunkStr := fmt.Sprintf("data: %s\n\n", string(msg.Data))
				fmt.Fprint(w, chunkStr)
				w.Flush()

			case <-time.After(5 * time.Minute): // Timeout de segurança
				log.Printf("⚠️ SSE Timeout para task %s", taskID)
				return
			}
		}
	})

	return nil
}
