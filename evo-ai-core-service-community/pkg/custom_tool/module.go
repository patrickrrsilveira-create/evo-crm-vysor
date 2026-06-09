package customTool

import (
	"evo-ai-core-service/pkg/custom_tool/handler"
	"evo-ai-core-service/pkg/custom_tool/repository"
	"evo-ai-core-service/pkg/custom_tool/service"

	"gorm.io/gorm"
)

type Module struct {
	Handler handler.CustomToolHandler
	Service service.CustomToolService
	Repo    repository.CustomToolRepository
}

func New(db *gorm.DB) *Module {
	r := repository.NewCustomToolRepository(db)
	s := service.NewCustomToolService(r)
	h := handler.NewCustomToolHandler(s)

	return &Module{
		Handler: h,
		Service: s,
		Repo:    r,
	}
}
