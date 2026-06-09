package customMcpServer

import (
	"evo-ai-core-service/internal/config"
	"evo-ai-core-service/pkg/custom_mcp_server/handler"
	"evo-ai-core-service/pkg/custom_mcp_server/repository"
	"evo-ai-core-service/pkg/custom_mcp_server/service"

	"gorm.io/gorm"
)

type Module struct {
	Handler handler.CustomMcpServerHandler
	Service service.CustomMcpServerService
	Repo    repository.CustomMcpServerRepository
}

func New(
	db *gorm.DB,
	cfgAIProcessorService *config.AIProcessorServiceConfig,
) *Module {
	r := repository.NewCustomMcpServerRepository(db)
	s := service.NewCustomMcpServerService(r, cfgAIProcessorService)
	h := handler.NewCustomMcpServerHandler(s)

	return &Module{
		Handler: h,
		Service: s,
		Repo:    r,
	}
}
