package repository

import (
	"context"
	"evo-campaign-engine/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AudienceRepo struct {
	db *gorm.DB
}

func NewAudienceRepo(db *gorm.DB) *AudienceRepo {
	return &AudienceRepo{db: db}
}

func (r *AudienceRepo) SaveAudience(ctx context.Context, a *domain.CampaignAudience) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AudienceRepo) SaveMembers(ctx context.Context, members []domain.AudienceMember) error {
	return r.db.WithContext(ctx).CreateInBatches(&members, 500).Error
}

func (r *AudienceRepo) GetByCampaign(ctx context.Context, campaignID uuid.UUID) (*domain.CampaignAudience, error) {
	var a domain.CampaignAudience
	err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AudienceRepo) ListMembers(ctx context.Context, campaignID uuid.UUID) ([]domain.AudienceMember, error) {
	var members []domain.AudienceMember
	err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Find(&members).Error
	return members, err
}

func (r *AudienceRepo) DeleteByCampaign(ctx context.Context, campaignID uuid.UUID) error {
	r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Delete(&domain.AudienceMember{})
	return r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Delete(&domain.CampaignAudience{}).Error
}
