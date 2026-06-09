package config

import (
	"encoding/json"
	"fmt"

	"evo-ai-core-service/internal/utils/stringutils"
	"evo-ai-core-service/pkg/agent/model"
)

func UpdateConfig(agent *model.Agent, updates map[string]interface{}) error {
	config := stringutils.JSONToInterfaceMap(agent.Config)
	if config == nil {
		config = make(map[string]interface{})
	}

	for k, v := range updates {
		config[k] = v
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	agent.Config = string(configJSON)
	return nil
}

func EnsureAPIKey(agent *model.Agent, generateKey func() string) error {
	config := stringutils.JSONToInterfaceMap(agent.Config)
	if config == nil {
		config = make(map[string]interface{})
	}

	if _, exists := config["api_key"]; !exists || config["api_key"] == "" {
		config["api_key"] = generateKey()
	}

	return UpdateConfig(agent, config)
}

func ConvertUUIDToStr(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		switch val := v.(type) {
		case map[string]interface{}:
			result[k] = ConvertUUIDToStr(val)
		case []interface{}:
			newSlice := make([]interface{}, len(val))
			for i, item := range val {
				if m, ok := item.(map[string]interface{}); ok {
					newSlice[i] = ConvertUUIDToStr(m)
				} else {
					newSlice[i] = item
				}
			}
			result[k] = newSlice
		default:
			result[k] = v
		}
	}
	return result
}
