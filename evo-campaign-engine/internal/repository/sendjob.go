package repository

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SendJobRepo struct {
	db *gorm.DB
}

func NewSendJobRepo(db *gorm.DB) *SendJobRepo {
	return &SendJobRepo{db: db}
}

func (r *SendJobRepo) BulkCreate(ctx context.Context, jobs []domain.SendJob) error {
	return r.db.WithContext(ctx).CreateInBatches(&jobs, 500).Error
}

func (r *SendJobRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.SendJob, error) {
	var j domain.SendJob
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&j).Error
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *SendJobRepo) UpdateState(ctx context.Context, id uuid.UUID, state domain.JobState, updates map[string]interface{}) error {
	if updates == nil {
		updates = make(map[string]interface{})
	}
	updates["state"] = state
	updates["updated_at"] = time.Now().UTC()
	return r.db.WithContext(ctx).Model(&domain.SendJob{}).Where("id = ?", id).Updates(updates).Error
}

func (r *SendJobRepo) FetchDueJobs(ctx context.Context, limit int) ([]domain.SendJob, error) {
	var jobs []domain.SendJob
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).
		Where("state = ? AND (scheduled_at IS NULL OR scheduled_at <= ?)", domain.JobQueued, now).
		Order("scheduled_at ASC NULLS FIRST").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}

func (r *SendJobRepo) CountByState(ctx context.Context, campaignID uuid.UUID) (map[domain.JobState]int, error) {
	type row struct {
		State string
		Count int
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Model(&domain.SendJob{}).
		Select("state, count(*) as count").
		Where("campaign_id = ?", campaignID).
		Group("state").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[domain.JobState]int, len(rows))
	for _, row := range rows {
		m[domain.JobState(row.State)] = row.Count
	}
	return m, nil
}

func (r *SendJobRepo) AssignSender(ctx context.Context, jobID, senderID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.SendJob{}).
		Where("id = ?", jobID).
		Update("sender_instance_id", senderID).Error
}

func (r *SendJobRepo) SaveDeliveryEvent(ctx context.Context, ev *domain.DeliveryEvent) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *SendJobRepo) Reschedule(ctx context.Context, jobID uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&domain.SendJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"scheduled_at": at,
			"state":        domain.JobScheduled,
			"updated_at":   time.Now().UTC(),
		}).Error
}
