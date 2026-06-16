package tools

import (
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

// GetFinanceTools retorna a lista de ferramentas que o FinanceAgent pode usar
func GetFinanceTools() []models.LLMTool {
	return []models.LLMTool{
		{
			Type: "function",
			Function: models.LLMFunctionDef{
				Name:        "generate_payment_link",
				Description: "Gera um link de pagamento (Stripe ou Mercado Pago) para cobrar o usuário final.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"provider": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"stripe", "mercadopago"},
							"description": "O provedor de pagamento escolhido.",
						},
						"amount": map[string]interface{}{
							"type":        "number",
							"description": "O valor da cobrança na moeda local (ex: 100.50).",
						},
						"currency": map[string]interface{}{
							"type":        "string",
							"description": "A moeda (ex: BRL, USD, EUR).",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Descrição do produto ou serviço.",
						},
					},
					"required": []string{"provider", "amount", "currency", "description"},
				},
			},
		},
	}
}


