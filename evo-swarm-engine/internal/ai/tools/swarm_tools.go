package tools

import (
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/swarm/registry"
)

// GetSwarmTools retorna ferramentas universais do sistema Swarm (como Delegação/Handoff),
// injetando dinamicamente a lista de agentes disponíveis baseada no Registry.
func GetSwarmTools(capabilities []registry.Capability, currentAgentID string) []models.LLMTool {
	var targetEnum []string

	for _, cap := range capabilities {
		if cap.AgentID != currentAgentID {
			targetEnum = append(targetEnum, cap.AgentID)
		}
	}

	// Se não houver outros agentes no Swarm, não precisamos injetar a tool de delegação.
	if len(targetEnum) == 0 {
		return []models.LLMTool{}
	}

	return []models.LLMTool{
		{
			Type: "function",
			Function: models.LLMFunctionDef{
				Name:        "delegate_to_agent",
				Description: "Delega a conversa atual e a intenção do usuário para um sub-agente especializado.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"target_agent": map[string]interface{}{
							"type":        "string",
							"description": "O ID do agente especialista que deve assumir a conversa.",
							"enum":        targetEnum,
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Motivo da delegação ou instruções de contexto sobre o que o sub-agente deve resolver para o usuário.",
						},
					},
					"required": []string{"target_agent", "reason"},
				},
			},
		},
	}
}
