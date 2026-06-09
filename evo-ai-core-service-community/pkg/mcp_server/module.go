package folderShare

import (
	"evo-ai-core-service/pkg/mcp_server/handler"
	"evo-ai-core-service/pkg/mcp_server/repository"
	"evo-ai-core-service/pkg/mcp_server/service"

	"gorm.io/gorm"
)

type Module struct {
	Handler handler.McpServerHandler
	Service service.McpServerService
	Repo    repository.McpServerRepository
}

func New(
	db *gorm.DB,
) *Module {
	r := repository.NewMcpServerRepository(db)
	s := service.NewMcpServerService(r)
	h := handler.NewMcpServerHandler(s)

	return &Module{
		Handler: h,
		Service: s,
		Repo:    r,
	}
}
