package model

import "evo-ai-core-service/internal/infra/postgres"

var MCPServerErrors = []postgres.CustomErrorMessage{
	{
		Code:    string(postgres.ERR_DUPLICATE_KEY_VIOLATION),
		Message: "MCP server with this URL already exists",
	},
	{
		Code:    string(postgres.ERR_RECORD_NOT_FOUND),
		Message: "MCP server not found",
	},
	{
		Code:    string(postgres.ERR_FOREIGN_KEY_VIOLATION),
		Message: "Invalid server configuration or reference",
	},
}
