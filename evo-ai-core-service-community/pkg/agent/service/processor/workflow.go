package processor

import (
	"context"

	"evo-ai-core-service/internal/utils/stringutils"
	"evo-ai-core-service/pkg/agent/model"
	"evo-ai-core-service/pkg/agent/service/config"
)

type WorkflowProcessor struct {
	generateAPIKey func() string
}

func NewWorkflowProcessor(generateAPIKey func() string) BaseProcessor {
	return &WorkflowProcessor{generateAPIKey: generateAPIKey}
}

func (p *WorkflowProcessor) Create(ctx context.Context, agent *model.Agent) error {
	agentConfig := stringutils.JSONToInterfaceMap(agent.Config)
	if agentConfig == nil {
		agentConfig = make(map[string]interface{})
	}

	agentConfig["api_key"] = p.generateAPIKey()

	return config.UpdateConfig(agent, agentConfig)
}

func (p *WorkflowProcessor) Update(ctx context.Context, current, request *model.Agent) error {
	if request.Type == model.AgentTypeWorkflow || request.Config != current.Config {
		updatedConfig := stringutils.JSONToInterfaceMap(request.Config)
		if updatedConfig == nil {
			updatedConfig = stringutils.JSONToInterfaceMap(current.Config)
		}

		currentConfig := stringutils.JSONToInterfaceMap(current.Config)
		if apiKey, exists := currentConfig["api_key"]; exists && apiKey != "" {
			updatedConfig["api_key"] = apiKey
		} else {
			updatedConfig["api_key"] = p.generateAPIKey()
		}

		return config.UpdateConfig(request, updatedConfig)
	}

	return nil
}
