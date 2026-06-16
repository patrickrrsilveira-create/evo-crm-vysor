package controllers

import (
	"strconv"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database/repositories"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AgentController gerencia as requisições HTTP relacionadas aos Agentes
type AgentController struct {
	repo repositories.AgentRepository
}

// NewAgentController cria uma nova instância de AgentController
func NewAgentController(repo repositories.AgentRepository) *AgentController {
	return &AgentController{repo: repo}
}

// Create lida com a criação de um novo agente
func (c *AgentController) Create(ctx *fiber.Ctx) error {
	var agent models.Agent
	if err := ctx.BodyParser(&agent); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido: " + err.Error()})
	}

	if agent.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "O campo 'name' é obrigatório"})
	}
	if agent.Type == "" {
		agent.Type = "llm" // Default type
	}

	if err := c.repo.Create(ctx.Context(), &agent); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao criar agente: " + err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Agente criado com sucesso",
		"data":    agent,
	})
}

// GetByID busca um agente pelo ID
func (c *AgentController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de agente inválido"})
	}

	agent, err := c.repo.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao buscar agente: " + err.Error()})
	}
	if agent == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Agente não encontrado"})
	}

	return ctx.JSON(fiber.Map{
		"message": "Agente recuperado com sucesso",
		"data":    agent,
	})
}

// List lista os agentes com paginação
func (c *AgentController) List(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	agents, total, err := c.repo.List(ctx.Context(), page, pageSize)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao listar agentes: " + err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"message":   "Agentes recuperados com sucesso",
		"data":      agents,
		"page":      page,
		"pageSize":  pageSize,
		"total":     total,
	})
}

// Update atualiza um agente existente
func (c *AgentController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de agente inválido"})
	}

	// Busca agente existente
	existingAgent, err := c.repo.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao buscar agente para atualização"})
	}
	if existingAgent == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Agente não encontrado"})
	}

	// Parse body
	var updateData models.Agent
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	// Update fields selectively (GORM's Save updates all non-zero fields if using struct, but Save with primary key updates everything)
	// For a true REST PUT, replace the object. For PATCH, should use Updates(). We'll just update fields explicitly for safety.
	existingAgent.Name = updateData.Name
	existingAgent.Description = updateData.Description
	existingAgent.Type = updateData.Type
	existingAgent.Model = updateData.Model
	existingAgent.Instruction = updateData.Instruction
	existingAgent.Role = updateData.Role
	existingAgent.Goal = updateData.Goal
	existingAgent.Config = updateData.Config
	if updateData.FolderID != nil {
		existingAgent.FolderID = updateData.FolderID
	}
	if updateData.APIKeyID != nil {
		existingAgent.APIKeyID = updateData.APIKeyID
	}

	if err := c.repo.Update(ctx.Context(), existingAgent); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao atualizar agente: " + err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"message": "Agente atualizado com sucesso",
		"data":    existingAgent,
	})
}

// Delete deleta um agente
func (c *AgentController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de agente inválido"})
	}

	if err := c.repo.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao deletar agente: " + err.Error()})
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
