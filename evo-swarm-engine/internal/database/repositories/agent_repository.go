package repositories

import (
	"context"
	"errors"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentRepository define a interface de acesso a dados para Agentes
type AgentRepository interface {
	Create(ctx context.Context, agent *models.Agent) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Agent, error)
	List(ctx context.Context, page, pageSize int) ([]models.Agent, int64, error)
	Update(ctx context.Context, agent *models.Agent) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type agentRepository struct {
	db *gorm.DB
}

// NewAgentRepository cria uma nova instância de AgentRepository
func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Create(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

func (r *agentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.WithContext(ctx).First(&agent, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ou um erro específico de não encontrado
		}
		return nil, err
	}
	return &agent, nil
}

func (r *agentRepository) List(ctx context.Context, page, pageSize int) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.Agent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&agents).Error
	if err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

func (r *agentRepository) Update(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

func (r *agentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Agent{}, "id = ?", id).Error
}
