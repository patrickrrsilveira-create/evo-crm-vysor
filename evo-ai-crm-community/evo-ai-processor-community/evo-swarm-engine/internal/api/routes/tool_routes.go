package routes

import (
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/api/controllers"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterToolRoutes registra as rotas REST para gerenciamento de Ferramentas Customizadas
func RegisterToolRoutes(app *fiber.App, db *gorm.DB) {
	repo := repositories.NewToolRepository(db)
	controller := controllers.NewToolController(repo)

	// Grupo de rotas base /api/v1/tools
	toolsGroup := app.Group("/api/v1/tools")

	toolsGroup.Post("/", controller.Create)
	toolsGroup.Get("/", controller.List)
	toolsGroup.Get("/:id", controller.GetByID)
	toolsGroup.Put("/:id", controller.Update)
	toolsGroup.Delete("/:id", controller.Delete)
}
