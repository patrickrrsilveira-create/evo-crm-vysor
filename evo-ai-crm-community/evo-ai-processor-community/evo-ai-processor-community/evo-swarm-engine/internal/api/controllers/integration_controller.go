package controllers

import (
	"encoding/json"
	"strings"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IntegrationController gerencia o Hot Swap dos endpoints de Canais do painel legado (Vue.js -> Python)
type IntegrationController struct {
	DB *gorm.DB
}

func NewIntegrationController(db *gorm.DB) *IntegrationController {
	return &IntegrationController{DB: db}
}

// sanitizeConfig remove senhas, tokens e secrets do JSON antes de mandar pro front-end
func sanitizeConfig(raw datatypes.JSON) map[string]interface{} {
	config := make(map[string]interface{})
	if err := json.Unmarshal(raw, &config); err != nil {
		return config
	}

	sensitiveFields := []string{
		"access_token", "client_id", "client_secret", "refresh_token",
		"token", "code_verifier", "password", "api_key",
	}

	for _, field := range sensitiveFields {
		delete(config, field)
	}

	// Remove também qualquer campo que comece com sk_ ou pk_ (ex: Stripe)
	for key, val := range config {
		if strVal, ok := val.(string); ok {
			if strings.HasPrefix(strVal, "sk_") || strings.HasPrefix(strVal, "pk_") || strings.HasPrefix(strVal, "rk_") {
				delete(config, key)
			}
		}
	}

	return config
}

// GetIntegrations lida com GET /agents/:id/integrations
func (c *IntegrationController) GetIntegrations(ctx *fiber.Ctx) error {
	agentIDStr := ctx.Params("id")
	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de Agente inválido"})
	}

	var integrations []models.AgentIntegration
	if err := c.DB.Where("agent_id = ?", agentID).Find(&integrations).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao buscar integrações"})
	}

	configs := make(map[string]interface{})
	credentialsConfigured := make(map[string]bool)

	for _, integ := range integrations {
		provider := integ.Provider
		
		var rawConfig map[string]interface{}
		json.Unmarshal(integ.Config, &rawConfig)
		
		// Determina se a integração está conectada baseada na existência de tokens antes de sanitizar
		connected := false
		if rawConfig["access_token"] != nil || rawConfig["token"] != nil || rawConfig["refresh_token"] != nil || rawConfig["api_key"] != nil {
			connected = true
		}

		if strings.HasSuffix(provider, "_credentials") {
			credentialsConfigured[provider] = connected
		} else {
			sanitized := sanitizeConfig(integ.Config)
			sanitized["connected"] = connected
			configs[provider] = sanitized
		}
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "All configurations retrieved successfully",
		"data": fiber.Map{
			"configs":                configs,
			"credentials_configured": credentialsConfigured,
		},
	})
}

// UpsertIntegration lida com POST /agents/:id/integrations
func (c *IntegrationController) UpsertIntegration(ctx *fiber.Ctx) error {
	agentIDStr := ctx.Params("id")
	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de Agente inválido"})
	}

	var body struct {
		Provider string                 `json:"provider"`
		Config   map[string]interface{} `json:"config"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido"})
	}

	normalizedProvider := strings.ReplaceAll(body.Provider, "-", "_")
	
	configJSON, _ := json.Marshal(body.Config)

	integ := models.AgentIntegration{
		AgentID:  agentID,
		Provider: normalizedProvider,
		Config:   configJSON,
		IsActive: true,
	}

	// Faz um Upsert nativo usando GORM
	err = c.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "agent_id"}, {Name: "provider"}},
		DoUpdates: clause.AssignmentColumns([]string{"config", "updated_at", "is_active"}),
	}).Create(&integ).Error

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao salvar integração"})
	}

	// Sanitiza a resposta antes de devolver para o frontend
	sanitizedResp := sanitizeConfig(configJSON)
	sanitizedResp["connected"] = true

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Integration " + normalizedProvider + " saved successfully",
		"data": fiber.Map{
			"provider":  normalizedProvider,
			"config":    sanitizedResp,
			"connected": true,
		},
	})
}

// DeleteIntegration lida com DELETE /agents/:id/integrations/:provider
func (c *IntegrationController) DeleteIntegration(ctx *fiber.Ctx) error {
	agentIDStr := ctx.Params("id")
	provider := ctx.Params("provider")

	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de Agente inválido"})
	}

	normalizedProvider := strings.ReplaceAll(provider, "-", "_")

	// Hard delete (mantendo compatibilidade com Python API)
	if err := c.DB.Where("agent_id = ? AND provider = ?", agentID, normalizedProvider).Delete(&models.AgentIntegration{}).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao excluir integração"})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Integration " + normalizedProvider + " deleted successfully",
		"data": fiber.Map{
			"provider": normalizedProvider,
			"deleted":  true,
		},
	})
}
