package validator

import (
	"context"
	"fmt"
	"log"

	"evo-ai-core-service/pkg/agent/model"

	"github.com/google/uuid"
)

type AgentValidator struct {
	getAgent func(ctx context.Context, id uuid.UUID) (*model.Agent, error)
}

func NewAgentValidator(getAgent func(ctx context.Context, id uuid.UUID) (*model.Agent, error)) *AgentValidator {
	return &AgentValidator{getAgent: getAgent}
}

func (v *AgentValidator) ValidateSubAgents(ctx context.Context, subAgents interface{}) error {
	subAgentsList, ok := subAgents.([]interface{})
	if !ok {
		return fmt.Errorf("sub_agents must be a list")
	}

	for _, sa := range subAgentsList {
		saID, ok := sa.(string)
		if !ok {
			return fmt.Errorf("invalid sub-agent ID format")
		}

		subAgentID, err := uuid.Parse(saID)
		if err != nil {
			return fmt.Errorf("invalid sub-agent ID: %v", err)
		}

		_, err = v.getAgent(ctx, subAgentID)
		if err != nil {
			return fmt.Errorf("sub-agent not found: %s", saID)
		}
	}
	return nil
}

func (v *AgentValidator) ValidateFlowConfig(ctx context.Context, config map[string]interface{}, agentType string) error {
	subAgents, ok := config["sub_agents"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid configuration: sub_agents must be a list for %s agents", agentType)
	}

	if len(subAgents) == 0 {
		return fmt.Errorf("invalid configuration: %s agents must have at least one sub-agent", agentType)
	}

	return v.ValidateSubAgents(ctx, subAgents)
}

func (v *AgentValidator) ValidateTaskConfig(ctx context.Context, config map[string]interface{}) error {
	tasks, ok := config["tasks"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid configuration: tasks is required")
	}

	if len(tasks) == 0 {
		return fmt.Errorf("invalid configuration: tasks cannot be empty")
	}

	for _, task := range tasks {
		if err := v.validateTask(ctx, task); err != nil {
			return err
		}
	}

	return nil
}

func (v *AgentValidator) validateTask(ctx context.Context, task interface{}) error {
	taskMap, ok := task.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid task configuration")
	}

	log.Println(taskMap)

	agentID, ok := taskMap["agent_id"].(string)
	if !ok {
		return fmt.Errorf("each task must have an agent_id")
	}

	taskAgentID, err := uuid.Parse(agentID)
	if err != nil {
		return fmt.Errorf("invalid agent_id format: %v", err)
	}

	_, err = v.getAgent(ctx, taskAgentID)
	if err != nil {
		return fmt.Errorf("agent not found for task: %s", agentID)
	}

	return nil
}
