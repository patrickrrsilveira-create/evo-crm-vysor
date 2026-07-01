package repository

import (
	"context"
	"evo-campaign-engine/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SuppressionRepo struct {
	db *gorm.DB
}

func NewSuppressionRepo(db *gorm.DB) *SuppressionRepo {
	return &SuppressionRepo{db: db}
}

func (r *SuppressionRepo) Add(ctx context.Context, s *domain.Suppression) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(s).Error
}

func (r *SuppressionRepo) IsSuppressed(ctx context.Context, accountID, contactID int, channelType string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Suppression{}).
		Where("account_id = ? AND contact_id = ? AND channel_type = ?", accountID, contactID, channelType).
		Count(&count).Error
	return count > 0, err
}

func (r *SuppressionRepo) Remove(ctx context.Context, accountID, contactID int, channelType string) error {
	return r.db.WithContext(ctx).
		Where("account_id = ? AND contact_id = ? AND channel_type = ?", accountID, contactID, channelType).
		Delete(&domain.Suppression{}).Error
}
