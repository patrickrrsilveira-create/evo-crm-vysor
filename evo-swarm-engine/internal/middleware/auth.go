package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// EvoAuthMiddleware valida JWTs e API Keys nativamente via GORM para latência ultra-baixa.
func EvoAuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Rotas de sistema e webhooks não precisam de autenticação
		path := c.Path()
		if path == "/health" || path == "/healthz" || path == "/readyz" || strings.HasPrefix(path, "/webhooks/") {
			return c.Next()
		}

		// 1. Tentar ler X-API-Key (Agent Bots / Scripts)
		apiKeyStr := c.Get("X-API-Key")

		// 2. Se não tiver X-API-Key, ler Authorization Bearer
		if apiKeyStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKeyStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Se após as duas tentativas não houver token, rejeita.
		if apiKeyStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"code":    "ERR_UNAUTHORIZED",
				"message": "Nenhum token de autenticação (Bearer ou X-API-Key) foi fornecido.",
			})
		}

		// 3. Validação Rápida no Banco de Dados (substituindo chamada HTTP de rede)
		// Primeiro verifica na tabela APIKey (Evo Core nativo)
		var keyRow models.APIKey
		err := db.Where("key = ?", apiKeyStr).First(&keyRow).Error

		if err == nil {
			// 4. Se encontrou, preenchemos o Agent Context fortemente tipado
			agentCtx := AgentContext{
				AgentID:   "system-agent", // Em produção mapearemos para evo_core_agents
				AgentName: keyRow.Name,
				KeyID:     keyRow.ID.String(),
			}

			c.Locals("AgentContext", agentCtx)
			c.Locals("is_agent_bot", true)

			return c.Next()
		}

		// Se não encontrou na tabela APIKey, verifica se é um Agent com API Key no config (padrão Python legado / Chatwoot / evo-bot-runtime)
		type Agent struct {
			ID   string
			Name string
		}
		var agentRow Agent
		// Query em jsonb para achar o api_key dentro do config
		errAgent := db.Table("agents").Select("id, name").Where("config->>'api_key' = ?", apiKeyStr).First(&agentRow).Error

		if errAgent == nil {
			agentCtx := AgentContext{
				AgentID:   agentRow.ID,
				AgentName: agentRow.Name,
				KeyID:     apiKeyStr,
			}

			c.Locals("AgentContext", agentCtx)
			c.Locals("is_agent_bot", true)

			return c.Next()
		}

		// Não encontrou em lugar nenhum no banco local. Fallback para validação remota (User Auth / Evo Auth).
		authBaseUrl := os.Getenv("EVO_AUTH_BASE_URL")
		if authBaseUrl == "" {
			authBaseUrl = "http://evo_auth:3001" // default docker-compose service
		}

		// Tenta validar no Evo Auth Service
		req, reqErr := http.NewRequest("POST", authBaseUrl+"/api/v1/auth/validate", nil)
		if reqErr == nil {
			req.Header.Set("Authorization", "Bearer "+apiKeyStr)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			client := &http.Client{Timeout: 5 * time.Second}
			resp, respErr := client.Do(req)

			if respErr == nil && resp.StatusCode == 200 {
				// Token validado com sucesso pelo Evo Auth!
				c.Locals("is_user_auth", true)
				if resp != nil {
					resp.Body.Close()
				}
				return c.Next()
			}
			if resp != nil {
				resp.Body.Close()
			}
		}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"code":    "ERR_INVALID_TOKEN",
				"message": "Token de acesso (API Key ou Bearer) inválido ou expirado.",
			})
		}

		// 4. Se encontrou, preenchemos o Agent Context fortemente tipado
		agentCtx := AgentContext{
			AgentID:   "system-agent", // Em produção mapearemos para evo_core_agents
			AgentName: keyRow.Name,
			KeyID:     keyRow.ID.String(),
		}

		c.Locals("AgentContext", agentCtx)
		c.Locals("is_agent_bot", true)

		return c.Next()
	}
}
