package repository

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SenderRepo struct {
	db *gorm.DB
}

func NewSenderRepo(db *gorm.DB) *SenderRepo {
	return &SenderRepo{db: db}
}

func (r *SenderRepo) Upsert(ctx context.Context, s *domain.SenderInstance) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *SenderRepo) GetByInboxIDs(ctx context.Context, inboxIDs []int) ([]domain.SenderInstance, error) {
	var senders []domain.SenderInstance
	err := r.db.WithContext(ctx).
		Where("inbox_id IN ? AND status = ?", inboxIDs, domain.SenderActive).
		Order("sent_today ASC, health_score DESC").
		Find(&senders).Error
	return senders, err
}

func (r *SenderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.SenderInstance, error) {
	var s domain.SenderInstance
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SenderRepo) IncrementSentToday(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&domain.SenderInstance{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"sent_today":  gorm.Expr("sent_today + 1"),
			"last_sent_at": now,
			"updated_at":   now,
		}).Error
}

func (r *SenderRepo) ResetDailyCounters(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&domain.SenderInstance{}).
		Where("sent_today > 0").
		Update("sent_today", 0).Error
}

func (r *SenderRepo) UpdateHealth(ctx context.Context, id uuid.UUID, score int, status domain.SenderStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.SenderInstance{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"health_score": score,
			"status":       status,
			"updated_at":   time.Now().UTC(),
		}).Error
}

func (r *SenderRepo) AdvanceWarmup(ctx context.Context, id uuid.UUID, day int) error {
	return r.db.WithContext(ctx).
		Model(&domain.SenderInstance{}).
		Where("id = ?", id).
		Update("warmup_stage_day", day).Error
}
