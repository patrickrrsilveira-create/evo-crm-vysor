package controllers

import (
	"strconv"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database/repositories"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ToolController gerencia as requisições HTTP relacionadas às Ferramentas Customizadas
type ToolController struct {
	repo repositories.ToolRepository
}

// NewToolController cria uma nova instância de ToolController
func NewToolController(repo repositories.ToolRepository) *ToolController {
	return &ToolController{repo: repo}
}

// Create lida com a criação de uma nova ferramenta
func (c *ToolController) Create(ctx *fiber.Ctx) error {
	var tool models.CustomTool
	if err := ctx.BodyParser(&tool); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido: " + err.Error()})
	}

	if tool.Name == "" || tool.Method == "" || tool.Endpoint == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Os campos 'name', 'method' e 'endpoint' são obrigatórios"})
	}

	if err := c.repo.Create(ctx.Context(), &tool); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao criar ferramenta: " + err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ferramenta criada com sucesso",
		"data":    tool,
	})
}

// GetByID busca uma ferramenta pelo ID
func (c *ToolController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de ferramenta inválido"})
	}

	tool, err := c.repo.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao buscar ferramenta: " + err.Error()})
	}
	if tool == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ferramenta não encontrada"})
	}

	return ctx.JSON(fiber.Map{
		"message": "Ferramenta recuperada com sucesso",
		"data":    tool,
	})
}

// List lista as ferramentas com paginação
func (c *ToolController) List(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	tools, total, err := c.repo.List(ctx.Context(), page, pageSize)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao listar ferramentas: " + err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"message":   "Ferramentas recuperadas com sucesso",
		"data":      tools,
		"page":      page,
		"pageSize":  pageSize,
		"total":     total,
	})
}

// Update atualiza uma ferramenta existente
func (c *ToolController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de ferramenta inválido"})
	}

	existingTool, err := c.repo.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao buscar ferramenta para atualização"})
	}
	if existingTool == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ferramenta não encontrada"})
	}

	var updateData models.CustomTool
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	// Campos atualizáveis
	existingTool.Name = updateData.Name
	existingTool.Description = updateData.Description
	existingTool.Method = updateData.Method
	existingTool.Endpoint = updateData.Endpoint
	existingTool.Headers = updateData.Headers
	existingTool.PathParams = updateData.PathParams
	existingTool.QueryParams = updateData.QueryParams
	existingTool.BodyParams = updateData.BodyParams
	existingTool.ErrorHandling = updateData.ErrorHandling
	existingTool.Values = updateData.Values
	existingTool.Tags = updateData.Tags
	existingTool.Examples = updateData.Examples
	existingTool.InputModes = updateData.InputModes
	existingTool.OutputModes = updateData.OutputModes

	if err := c.repo.Update(ctx.Context(), existingTool); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao atualizar ferramenta: " + err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"message": "Ferramenta atualizada com sucesso",
		"data":    existingTool,
	})
}

// Delete deleta uma ferramenta
func (c *ToolController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de ferramenta inválido"})
	}

	if err := c.repo.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao deletar ferramenta: " + err.Error()})
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
