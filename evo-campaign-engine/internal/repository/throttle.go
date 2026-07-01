package repository

import (
	"context"
	"evo-campaign-engine/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThrottleRepo struct {
	db *gorm.DB
}

func NewThrottleRepo(db *gorm.DB) *ThrottleRepo {
	return &ThrottleRepo{db: db}
}

func (r *ThrottleRepo) Create(ctx context.Context, p *domain.ThrottleProfile) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ThrottleRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ThrottleProfile, error) {
	var p domain.ThrottleProfile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ThrottleRepo) ListByChannel(ctx context.Context, channelType string) ([]domain.ThrottleProfile, error) {
	var profiles []domain.ThrottleProfile
	err := r.db.WithContext(ctx).Where("channel_type = ?", channelType).Find(&profiles).Error
	return profiles, err
}

func (r *ThrottleRepo) GetDefault(ctx context.Context, channelType string) (*domain.ThrottleProfile, error) {
	var p domain.ThrottleProfile
	err := r.db.WithContext(ctx).Where("channel_type = ?", channelType).Order("created_at ASC").First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ThrottleRepo) Update(ctx context.Context, p *domain.ThrottleProfile) error {
	return r.db.WithContext(ctx).Save(p).Error
}
