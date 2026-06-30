package routes

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterIntegrationRoutes registra os endpoints para consumir recursos das integrações configuradas
func RegisterIntegrationRoutes(app *fiber.App) {
	// Cria um subgrupo para integrations
	api := app.Group("/api/v1/integrations")

	// Endpoint genérico para listar recursos de um provedor (ex: Calendários, Bases do Notion)
	// Isso permite que o Frontend UI mostre dropdowns de seleção sem precisar conhecer a API de cada serviço.
	api.Get("/:provider/resources", func(c *fiber.Ctx) error {
		provider := c.Params("provider")

		// TODO: Implementar ResourceFetcher genérico que interage com o banco de dados
		// para recuperar o token do agent/user e bater na API do provider

		switch provider {
		case "google_calendar":
			return c.JSON(fiber.Map{
				"status": "success",
				"data": []fiber.Map{
					{"id": "primary", "name": "Calendário Principal"},
				},
			})
		case "notion":
			return c.JSON(fiber.Map{
				"status": "success",
				"data": []fiber.Map{
					{"id": "mock_db_1", "name": "Leads 2026"},
				},
			})
		default:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "provedor não suportado para listagem de recursos",
			})
		}
	})
}
