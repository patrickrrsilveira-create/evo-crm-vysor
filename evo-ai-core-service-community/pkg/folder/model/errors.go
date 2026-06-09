package model

import "evo-ai-core-service/internal/infra/postgres"

var FolderErrors = []postgres.CustomErrorMessage{
	{
		Code:    string(postgres.ERR_DUPLICATE_KEY_VIOLATION),
		Message: "Folder with this name already exists in this path",
	},
	{
		Code:    string(postgres.ERR_RECORD_NOT_FOUND),
		Message: "Folder not found",
	},
	{
		Code:    string(postgres.ERR_FOREIGN_KEY_VIOLATION),
		Message: "Cannot delete folder with associated resources",
	},
}
