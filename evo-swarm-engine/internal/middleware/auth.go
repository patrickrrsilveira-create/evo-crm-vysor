package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// EvoAuthMiddleware simula o validador JWT/Bearer do EvoAuth
func EvoAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Rotas de sistema e webhooks abertos não precisam de autenticação profunda
		path := c.Path()
		if path == "/health" || path == "/healthz" || path == "/readyz" || strings.HasPrefix(path, "/webhooks/") {
			return c.Next()
		}

		// Validação estrita de Header
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Acesso negado. Token de autenticação não fornecido ou inválido.",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Aqui entraria a validação real do Token via serviço gRPC ou chamada ao BD
		// Exemplo simulado de rejeição:
		if token == "token_expirado_ou_invalido" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Token expirado.",
			})
		}

		// Injetando dados do usuário no contexto
		c.Locals("userID", "user_123")

		return c.Next()
	}
}
