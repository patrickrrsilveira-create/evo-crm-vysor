package processor

import (
	"context"
	"fmt"
	"regexp"

	"evo-ai-core-service/internal/utils/stringutils"
	"evo-ai-core-service/pkg/agent/client/a2a"
	"evo-ai-core-service/pkg/agent/model"
	"evo-ai-core-service/pkg/agent/service/config"
)

type A2AProcessor struct {
	a2aClient a2a.Client
}

func NewA2AProcessor(a2aClient a2a.Client) BaseProcessor {
	return &A2AProcessor{a2aClient: a2aClient}
}

func (p *A2AProcessor) Create(ctx context.Context, agent *model.Agent) error {
	if agent.CardURL == "" {
		return fmt.Errorf("card_url is required for a2a type agents")
	}

	agentCard, err := p.a2aClient.FetchAgentCard(ctx, agent.CardURL)
	if err != nil {
		return err
	}

	return p.updateAgentFromCard(agent, agentCard, true)
}

func (p *A2AProcessor) Update(ctx context.Context, current, request *model.Agent) error {
	if (request.Type == model.AgentTypeA2A || current.Type == model.AgentTypeA2A) && request.CardURL != "" {

		agentCard, err := p.a2aClient.FetchAgentCard(ctx, request.CardURL)
		if err != nil {
			return err
		}

		return p.updateAgentFromCard(request, agentCard, false)
	}
	return nil
}

func (p *A2AProcessor) updateAgentFromCard(agent *model.Agent, agentCard map[string]interface{}, isCreate bool) error {
	if isCreate || agent.Name == "" {
		cardName, _ := agentCard["name"].(string)
		if cardName == "" {
			cardName = "Unknown Agent"
		}
		agent.Name = sanitizeName(cardName)
	}

	if isCreate || agent.Description == "" {
		if desc, ok := agentCard["description"].(string); ok {
			agent.Description = desc
		}
	}

	agentConfig := stringutils.JSONToInterfaceMap(agent.Config)
	if agentConfig == nil {
		agentConfig = make(map[string]interface{})
	}

	agentConfig["agent_card"] = agentCard

	return config.UpdateConfig(agent, agentConfig)
}

func sanitizeName(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return reg.ReplaceAllString(name, "_")
}
