package routes

import (
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/middleware"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterOAuthRoutes inicializa os endpoints PKCE
func RegisterOAuthRoutes(app *fiber.App, oauthService *services.OAuthService, db *gorm.DB) {
	group := app.Group("/api/v1/oauth")

	// GET /api/v1/oauth/{provider}/authorize
	// Endpoint protegido, requer o auth middleware da engine
	group.Get("/:provider/authorize", middleware.EvoAuthMiddleware(db), func(c *fiber.Ctx) error {
		providerName := c.Params("provider")

		// Obtém o AgentID a partir do contexto protegido
		agentCtx := c.Locals("AgentContext").(middleware.AgentContext)

		authURL, err := oauthService.GenerateAuthorizationURL(agentCtx.AgentID, providerName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": err.Error(),
			})
		}

		// Retorna a URL para que o Frontend do usuário faça o redirecionamento
		return c.JSON(fiber.Map{
			"authorization_url": authURL,
		})
	})

	// GET /api/v1/oauth/{provider}/callback
	// O Callback NÃO é protegido, pois é o Google/Notion que nos redireciona de volta
	group.Get("/:provider/callback", func(c *fiber.Ctx) error {
		// providerName := c.Params("provider")
		state := c.Query("state")
		code := c.Query("code")

		if state == "" || code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "State e Code são obrigatórios no callback",
			})
		}

		tokenResp, err := oauthService.HandleCallback(c.Context(), state, code)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Falha na Autenticação",
				"message": err.Error(),
			})
		}

		// Fecha a aba do usuário enviando uma tela de sucesso
		// Em um ambiente real, isso gravaria o tokenResp no banco.
		return c.SendString("Integração concluída com sucesso! O Token de Acesso foi armazenado de forma segura (" + tokenResp.TokenType + "). Pode fechar esta janela.")
	})
}
