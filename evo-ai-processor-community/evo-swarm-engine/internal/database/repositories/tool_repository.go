package repositories

import (
	"context"
	"errors"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ToolRepository define a interface de acesso a dados para Ferramentas Customizadas
type ToolRepository interface {
	Create(ctx context.Context, tool *models.CustomTool) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.CustomTool, error)
	List(ctx context.Context, page, pageSize int) ([]models.CustomTool, int64, error)
	Update(ctx context.Context, tool *models.CustomTool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type toolRepository struct {
	db *gorm.DB
}

// NewToolRepository cria uma nova instância de ToolRepository
func NewToolRepository(db *gorm.DB) ToolRepository {
	return &toolRepository{db: db}
}

func (r *toolRepository) Create(ctx context.Context, tool *models.CustomTool) error {
	return r.db.WithContext(ctx).Create(tool).Error
}

func (r *toolRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CustomTool, error) {
	var tool models.CustomTool
	err := r.db.WithContext(ctx).First(&tool, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ou um erro específico
		}
		return nil, err
	}
	return &tool, nil
}

func (r *toolRepository) List(ctx context.Context, page, pageSize int) ([]models.CustomTool, int64, error) {
	var tools []models.CustomTool
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.CustomTool{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&tools).Error
	if err != nil {
		return nil, 0, err
	}

	return tools, total, nil
}

func (r *toolRepository) Update(ctx context.Context, tool *models.CustomTool) error {
	return r.db.WithContext(ctx).Save(tool).Error
}

func (r *toolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.CustomTool{}, "id = ?", id).Error
}
