package middleware

import (
	"encoding/json"
	"fmt"
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

		// Extrair agent_id do path (ex: /api/v1/a2a/1234-5678-...)
		parts := strings.Split(path, "/")
		var agentIDStr string
		for i, part := range parts {
			if part == "a2a" || part == "chat" {
				if len(parts) > i+1 {
					agentIDStr = parts[i+1]
					break
				}
			}
		}

		// 1. Tentar ler X-API-Key (Agent Bots / Scripts)
		apiKeyStr := c.Get("X-API-Key")

		// 2. Se não tiver X-API-Key, ler Authorization Bearer
		if apiKeyStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKeyStr = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
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
		// Primeiro, tenta validar como API Key de sistema (evo_core_api_keys)
		var keyRow models.APIKey
		err := db.Where("key = ?", apiKeyStr).First(&keyRow).Error

		if err == nil {
			agentCtx := AgentContext{
				AgentID:   "system-agent",
				AgentName: keyRow.Name,
				KeyID:     keyRow.ID.String(),
			}

			c.Locals("AgentContext", agentCtx)
			c.Locals("is_agent_bot", true)

			return c.Next()
		}

		// 4. Para rotas A2A/Chat, valida contra o api_key do agente no evo_core_agents
		// O bot-runtime envia X-API-Key que deve corresponder ao agent.config.api_key
		if agentIDStr != "" {
			type Agent struct {
				ID     string
				Name   string
				Config string
			}
			var agentRow Agent
			errAgent := db.Table("evo_core_agents").Select("id, name, config::text as config").Where("id = ?", agentIDStr).First(&agentRow).Error

			if errAgent == nil && agentRow.Config != "" {
				var configMap map[string]interface{}
				if errParse := json.Unmarshal([]byte(agentRow.Config), &configMap); errParse == nil {
					if storedKey, ok := configMap["api_key"].(string); ok {
						if storedKey == apiKeyStr {
							agentCtx := AgentContext{
								AgentID:   agentRow.ID,
								AgentName: agentRow.Name,
								KeyID:     apiKeyStr,
							}

							c.Locals("AgentContext", agentCtx)
							c.Locals("is_agent_bot", true)

							return c.Next()
						}
						// API Key não confere - log para debug
						fmt.Printf("⚠️ [AuthDebug] Invalid API key for agent %s: provided=%s... stored=%s...\n", agentIDStr, apiKeyStr[:min(8, len(apiKeyStr))], storedKey[:min(8, len(storedKey))])
					} else {
						fmt.Printf("⚠️ [AuthDebug] Agent %s has no api_key in config\n", agentIDStr)
					}
				}
			} else if errAgent != nil {
				fmt.Printf("⚠️ [AuthDebug] Agent not found in DB for A2A call: %s (error: %v)\n", agentIDStr, errAgent)
			}
		} else {
			fmt.Printf("⚠️ [AuthDebug] Could not extract agent_id from path: %s\n", path)
		}

		// 5. Não encontrou em lugar nenhum no banco local. Fallback para validação remota.
		authBaseUrl := os.Getenv("EVO_AUTH_BASE_URL")
		if authBaseUrl == "" {
			authBaseUrl = "http://evo_auth:3001"
		}

		req, reqErr := http.NewRequest("POST", authBaseUrl+"/api/v1/auth/validate", nil)
		if reqErr == nil {
			req.Header.Set("Authorization", "Bearer "+apiKeyStr)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			client := &http.Client{Timeout: 5 * time.Second}
			resp, respErr := client.Do(req)

			if respErr == nil && resp.StatusCode == 200 {
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
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
