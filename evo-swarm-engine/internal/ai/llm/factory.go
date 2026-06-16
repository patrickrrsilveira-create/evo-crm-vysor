package llm

import (
	"fmt"
	"strings"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

// NewLLMProvider retorna o cliente correto baseado no nome do modelo e chave de API
func NewLLMProvider(modelName string, apiKey string) (models.LLMProvider, error) {
	lowerModel := strings.ToLower(modelName)

	if strings.HasPrefix(lowerModel, "gpt-") || strings.HasPrefix(lowerModel, "o1-") {
		return NewOpenAIClient(apiKey), nil
	}

	if strings.HasPrefix(lowerModel, "claude-") {
		return NewAnthropicClient(apiKey), nil
	}

	if strings.HasPrefix(lowerModel, "gemini-") {
		return NewGeminiClient(apiKey), nil
	}

	return nil, fmt.Errorf("provedor LLM não suportado para o modelo: %s", modelName)
}
