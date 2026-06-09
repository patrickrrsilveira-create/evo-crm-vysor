package model

import "evo-ai-core-service/internal/infra/postgres"

var FolderShareErrors = []postgres.CustomErrorMessage{
	{
		Code:    string(postgres.ERR_DUPLICATE_KEY_VIOLATION),
		Message: "Folder is already shared with this client",
	},
	{
		Code:    string(postgres.ERR_RECORD_NOT_FOUND),
		Message: "Shared folder not found",
	},
	{
		Code:    string(postgres.ERR_FOREIGN_KEY_VIOLATION),
		Message: "Invalid folder or client reference",
	},
}
