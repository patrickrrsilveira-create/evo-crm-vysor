package routes

import (
	"context"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/mcp"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// MCPDiscoverRequest representa a estrutura que o Evo CRM enviará
type MCPDiscoverRequest struct {
	URL string `json:"url"` // A URL do servidor MCP customizado (ex: ngrok, Railway, etc)
}

// RegisterMCPRoutes expõe os endpoints para gerenciar servidores Model Context Protocol (MCP) externos.
func RegisterMCPRoutes(app *fiber.App, db *gorm.DB) {
	group := app.Group("/api/v1/custom-mcp-servers")

	// GET /discover-tools descobre ferramentas em tempo real de forma dinâmica
	group.Post("/discover-tools", middleware.EvoAuthMiddleware(db), func(c *fiber.Ctx) error {
		var req MCPDiscoverRequest
		if err := c.BodyParser(&req); err != nil || req.URL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "URL do MCP Server é obrigatória",
			})
		}

		// Criação dinâmica e On-the-fly do Client MCP (abandonando o singleton hardcoded de main.go)
		dynamicClient := mcp.NewClient(req.URL)

		ctx, cancel := context.WithTimeout(context.Background(), 10*60*1000*1000*1000) // 10s
		defer cancel()

		if err := dynamicClient.Connect(ctx); err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "Bad Gateway",
				"message": "Falha ao conectar no servidor MCP remoto: " + err.Error(),
			})
		}

		// O cliente faz o fetch das ferramentas remotas
		tools, err := dynamicClient.DiscoverTools(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Falha ao descobrir ferramentas: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"tools":   tools,
			"message": "Ferramentas MCP descobertas com sucesso na URL.",
		})
	})
}
