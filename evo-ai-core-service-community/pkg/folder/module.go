package folder

import (
	"evo-ai-core-service/pkg/folder/handler"
	"evo-ai-core-service/pkg/folder/repository"
	"evo-ai-core-service/pkg/folder/service"

	"gorm.io/gorm"
)

type Module struct {
	Handler handler.FolderHandler
	Service service.FolderService
	Repo    repository.FolderRepository
}

func New(
	db *gorm.DB,
) *Module {
	r := repository.NewFolderRepository(db)
	s := service.NewFolderService(r)
	h := handler.NewFolderHandler(s)

	return &Module{
		Handler: h,
		Service: s,
		Repo:    r,
	}
}
