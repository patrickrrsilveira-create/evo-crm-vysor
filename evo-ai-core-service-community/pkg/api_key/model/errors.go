package model

import "evo-ai-core-service/internal/infra/postgres"

var APIKeyErrors = []postgres.CustomErrorMessage{
	{
		Code:    string(postgres.ERR_DUPLICATE_KEY_VIOLATION),
		Message: "API key with this name already exists",
	},
	{
		Code:    string(postgres.ERR_RECORD_NOT_FOUND),
		Message: "API key not found",
	},
	{
		Code:    string(postgres.ERR_FOREIGN_KEY_VIOLATION),
		Message: "Client not found or invalid client reference",
	},
}
