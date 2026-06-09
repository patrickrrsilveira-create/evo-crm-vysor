package processor

import (
	"context"
	"fmt"

	"evo-ai-core-service/internal/utils/stringutils"
	"evo-ai-core-service/pkg/agent/model"
	"evo-ai-core-service/pkg/agent/service/config"
)

// Supported external providers
var validExternalProviders = []string{
	"flowise", // Flowise AI chatflows
	"n8n",     // N8N workflows
	"typebot", // Typebot chatbots
	"dify",    // Dify AI agents
	"openai",  // OpenAI Assistants/Chat Completion
}

type ExternalProcessor struct {
	generateAPIKey func() string
}

func NewExternalProcessor(generateAPIKey func() string) BaseProcessor {
	return &ExternalProcessor{
		generateAPIKey: generateAPIKey,
	}
}

func (p *ExternalProcessor) Create(ctx context.Context, agent *model.Agent) error {
	agentConfig := stringutils.JSONToInterfaceMap(agent.Config)
	if agentConfig == nil {
		agentConfig = make(map[string]interface{})
	}

	// Validate that provider is present
	provider, ok := agentConfig["provider"].(string)
	if !ok || provider == "" {
		return fmt.Errorf("provider is required for external type agents")
	}

	// Validate provider value - must be one of the supported providers
	valid := false
	for _, vp := range validExternalProviders {
		if provider == vp {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid provider: %s. Must be one of: %v", provider, validExternalProviders)
	}

	// Generate API key
	agentConfig["api_key"] = p.generateAPIKey()

	return config.UpdateConfig(agent, agentConfig)
}

func (p *ExternalProcessor) Update(ctx context.Context, current, request *model.Agent) error {
	updatedConfig := stringutils.JSONToInterfaceMap(request.Config)
	if updatedConfig == nil {
		updatedConfig = stringutils.JSONToInterfaceMap(current.Config)
	}

	// Validate provider if present - must be one of the supported providers
	if provider, ok := updatedConfig["provider"].(string); ok && provider != "" {
		valid := false
		for _, vp := range validExternalProviders {
			if provider == vp {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid provider: %s. Must be one of: %v", provider, validExternalProviders)
		}
	}

	// Preserve API key from current config if exists
	currentConfig := stringutils.JSONToInterfaceMap(current.Config)
	if apiKey, exists := currentConfig["api_key"]; exists && apiKey != "" {
		updatedConfig["api_key"] = apiKey
	} else {
		updatedConfig["api_key"] = p.generateAPIKey()
	}

	return config.UpdateConfig(request, updatedConfig)
}
