package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

// MSTeamsAdapter gerencia o Microsoft Bot Framework (Teams)
type MSTeamsAdapter struct {
	EventBus     *evbus.EventBus
	AppID        string
	AppPassword  string
	accessToken  string
	tokenExpires time.Time
	mu           sync.RWMutex
}

func NewMSTeamsAdapter(bus *evbus.EventBus) *MSTeamsAdapter {
	return &MSTeamsAdapter{
		EventBus:    bus,
		AppID:       os.Getenv("MSTEAMS_APP_ID"),
		AppPassword: os.Getenv("MSTEAMS_APP_PASSWORD"),
	}
}

func (a *MSTeamsAdapter) Name() string {
	return "msteams"
}

func (a *MSTeamsAdapter) Start(ctx context.Context) error {
	log.Println("🟦 [MSTeamsAdapter] Inicializado (Aguardando Webhooks da Microsoft)")
	
	// Escuta respostas geradas pelos Agentes do Swarm
	_, err := a.EventBus.Conn.QueueSubscribe("outbound.message", "msteams_outbound_group", a.handleOutbound)
	return err
}

func (a *MSTeamsAdapter) handleOutbound(msg *nats.Msg) {
	var payload map[string]interface{}
	json.Unmarshal(msg.Data, &payload)

	source, _ := payload["source"].(string)
	if source != "msteams" {
		return // Não é pra gente
	}

	sender, _ := payload["sender"].(string)
	content, _ := payload["content"].(string)

	if sender != "" && content != "" {
		err := a.SendMessage(context.Background(), sender, content)
		if err != nil {
			log.Printf("❌ [MSTeamsAdapter] Erro ao enviar resposta: %v", err)
		}
	}
}

// getToken garante um token JWT de autenticação usando Client Credentials flow.
// Realiza cache do token na memória para performance extrema.
func (a *MSTeamsAdapter) getToken(ctx context.Context) (string, error) {
	a.mu.RLock()
	if a.accessToken != "" && time.Now().Before(a.tokenExpires) {
		token := a.accessToken
		a.mu.RUnlock()
		return token, nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	// Double check pattern
	if a.accessToken != "" && time.Now().Before(a.tokenExpires) {
		return a.accessToken, nil
	}

	authURL := "https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.AppID)
	data.Set("client_secret", a.AppPassword)
	data.Set("scope", "https://api.botframework.com/.default")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("falha ao obter token ms teams: %d %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Subtraímos 5 minutos do TTL para termos margem de segurança (jitter)
	a.accessToken = result.AccessToken
	a.tokenExpires = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)

	log.Println("🟦 [MSTeamsAdapter] Novo Token Microsoft Auth gerado em cache")
	return a.accessToken, nil
}

// SendMessage envia uma mensagem de volta ao Teams usando o Bot Connector REST API
// A variável 'to' deve conter o Conversation ID ou ServiceURL empacotado.
// Por simplificação do motor, vamos assumir que 'to' = {serviceUrl}|{conversationId}
func (a *MSTeamsAdapter) SendMessage(ctx context.Context, to string, content string) error {
	parts := strings.SplitN(to, "|", 2)
	if len(parts) != 2 {
		return fmt.Errorf("formato do destinatário inválido para teams, experado: serviceUrl|conversationId")
	}
	serviceURL := parts[0]
	conversationID := parts[1]

	token, err := a.getToken(ctx)
	if err != nil {
		return fmt.Errorf("msteams auth error: %w", err)
	}

	// Rota do Bot Framework: {serviceUrl}/v3/conversations/{conversationId}/activities
	endpoint := fmt.Sprintf("%s/v3/conversations/%s/activities", serviceURL, url.PathEscape(conversationID))

	payload := map[string]interface{}{
		"type": "message",
		"text": content,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao enviar msg ms teams (%d): %s", resp.StatusCode, string(body))
	}

	log.Printf("🟦 [MSTeamsAdapter] Mensagem enviada com sucesso para conversa %s", conversationID)
	return nil
}

// RegisterWebhookRoute registra a rota HTTP de recepção de atividades (Incoming Messages)
func (a *MSTeamsAdapter) RegisterWebhookRoute(app *fiber.App) {
	app.Post("/api/v1/webhooks/msteams", func(c *fiber.Ctx) error {
		// Eventos do Microsoft Bot Framework chegam como Array de JSON Objects
		// ou como um único Object chamado "Activity".

		var payload map[string]interface{}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString("Bad Request")
		}

		activityType, _ := payload["type"].(string)

		// Estamos interessados apenas em mensagens enviadas por humanos
		if activityType == "message" {
			text, _ := payload["text"].(string)
			serviceURL, _ := payload["serviceUrl"].(string)

			// Extrai a conversa para que a IA possa responder depois
			var conversationID string
			if conv, ok := payload["conversation"].(map[string]interface{}); ok {
				conversationID, _ = conv["id"].(string)
			}

			if text != "" && serviceURL != "" && conversationID != "" {
				// Formato do remetente para a IA conseguir responder: serviceUrl|conversationId
				senderRef := fmt.Sprintf("%s|%s", serviceURL, conversationID)

				eventData, _ := json.Marshal(map[string]interface{}{
					"source":  "msteams",
					"sender":  senderRef,
					"content": text,
				})

				log.Printf("📥 [MSTeamsAdapter] Nova mensagem recebida. Despachando no NATS...")
				a.EventBus.Publish(string(events.EventMessageReceived), eventData)
			}
		}

		return c.SendStatus(200)
	})
}
