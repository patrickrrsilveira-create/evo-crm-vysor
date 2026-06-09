package model

import "evo-ai-core-service/internal/infra/postgres"

var AgentErrors = []postgres.CustomErrorMessage{
	{
		Code:    string(postgres.ERR_DUPLICATE_KEY_VIOLATION),
		Message: "Agent with this name already exists for this client",
	},
	{
		Code:    string(postgres.ERR_RECORD_NOT_FOUND),
		Message: "Agent not found",
	},
	{
		Code:    string(postgres.ERR_FOREIGN_KEY_VIOLATION),
		Message: "Client not found or invalid client reference",
	},
}
