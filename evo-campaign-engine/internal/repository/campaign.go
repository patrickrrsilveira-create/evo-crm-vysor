package repository

import (
	"context"
	"evo-campaign-engine/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CampaignRepo struct {
	db *gorm.DB
}

func NewCampaignRepo(db *gorm.DB) *CampaignRepo {
	return &CampaignRepo{db: db}
}

func (r *CampaignRepo) Create(ctx context.Context, c *domain.Campaign) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CampaignRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	var c domain.Campaign
	err := r.db.WithContext(ctx).
		Preload("Channels").
		Preload("Variants").
		Preload("Audience").
		Where("id = ?", id).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CampaignRepo) ListByAccount(ctx context.Context, accountID, page, pageSize int) ([]domain.Campaign, int64, error) {
	var campaigns []domain.Campaign
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Campaign{}).Where("account_id = ?", accountID)
	q.Count(&total)

	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&campaigns).Error

	return campaigns, total, err
}

func (r *CampaignRepo) Update(ctx context.Context, c *domain.Campaign) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CampaignRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CampaignStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Campaign{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *CampaignRepo) IncrementCounter(ctx context.Context, id uuid.UUID, field string, delta int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Campaign{}).
		Where("id = ?", id).
		UpdateColumn(field, gorm.Expr(field+" + ?", delta)).Error
}

func (r *CampaignRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Campaign{}).Error
}

func (r *CampaignRepo) SaveChannels(ctx context.Context, channels []domain.CampaignChannel) error {
	return r.db.WithContext(ctx).Create(&channels).Error
}

func (r *CampaignRepo) DeleteChannels(ctx context.Context, campaignID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Delete(&domain.CampaignChannel{}).Error
}

func (r *CampaignRepo) GetChannels(ctx context.Context, campaignID uuid.UUID) ([]domain.CampaignChannel, error) {
	var channels []domain.CampaignChannel
	err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Find(&channels).Error
	return channels, err
}

func (r *CampaignRepo) SaveVariants(ctx context.Context, variants []domain.MessageVariant) error {
	return r.db.WithContext(ctx).Create(&variants).Error
}

func (r *CampaignRepo) GetVariants(ctx context.Context, campaignID uuid.UUID) ([]domain.MessageVariant, error) {
	var variants []domain.MessageVariant
	err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).Find(&variants).Error
	return variants, err
}
