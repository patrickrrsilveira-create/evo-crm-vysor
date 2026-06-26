package routes

import (
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/api/controllers"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterAgentRoutes registra as rotas REST para gerenciamento de Agentes
func RegisterAgentRoutes(app *fiber.App, db *gorm.DB) {
	// Inicializar dependências (Injeção de Dependência)
	repo := repositories.NewAgentRepository(db)
	controller := controllers.NewAgentController(repo)

	// Grupo de rotas base /api/v1/agents
	agentsGroup := app.Group("/api/v1/agents")

	agentsGroup.Post("/", controller.Create)
	agentsGroup.Get("/", controller.List)
	agentsGroup.Get("/:id", controller.GetByID)
	agentsGroup.Put("/:id", controller.Update)
	agentsGroup.Delete("/:id", controller.Delete)

	// Inicializa o Integration Controller
	integrationController := controllers.NewIntegrationController(db)
	
	// Hot Swap: Mesmos endpoints que a API legada do Python usava para a página de Canais do CRM
	agentsGroup.Get("/:id/integrations", integrationController.GetIntegrations)
	agentsGroup.Post("/:id/integrations", integrationController.UpsertIntegration)
	agentsGroup.Delete("/:id/integrations/:provider", integrationController.DeleteIntegration)
}
