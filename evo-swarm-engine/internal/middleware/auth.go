package middleware

import (
	"strings"

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
		var keyRow models.APIKey
		err := db.Where("key = ?", apiKeyStr).First(&keyRow).Error

		if err != nil {
			// Não encontrou na tabela APIKey. Aqui seria o fallback para validação de assinatura JWT (User Auth).
			// Para MVP/Paridade, rejeitamos.
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"code":    "ERR_INVALID_API_KEY",
				"message": "API Key inválida ou não encontrada.",
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
