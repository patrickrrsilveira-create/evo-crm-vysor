package tools

import (
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

// GetCRMTools retorna a lista de ferramentas nativas que a IA pode invocar para manipular o CRM
func GetCRMTools() []models.LLMTool {
	return []models.LLMTool{
		{
			Type: "function",
			Function: models.LLMFunctionDef{
				Name:        "transfer_to_human",
				Description: "Transfere a conversa atual para um atendente humano. Use esta ferramenta quando o usuário pedir para falar com um humano, estiver frustrado, ou quando a complexidade for maior do que você pode resolver.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Motivo da transferência (ex: Usuário solicitou falar com atendente)",
						},
						// Nota: Na versão Go, a IA não precisa adivinhar o conversation_id, ele será injetado pelo Contexto do AgentWorker.
					},
					"required": []string{"reason"},
				},
			},
		},
	}
}


