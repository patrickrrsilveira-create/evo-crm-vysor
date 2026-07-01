package throttle

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Engine struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Engine {
	return &Engine{rdb: rdb}
}

func (e *Engine) WarmupCap(profile *domain.ThrottleProfile, daysSinceStart int) int {
	if profile.WarmupStepDays <= 0 {
		return profile.DailyCapMax
	}
	steps := daysSinceStart / profile.WarmupStepDays
	cap := float64(profile.DailyCapStart) * math.Pow(profile.WarmupMultiplier, float64(steps))
	result := int(math.Floor(cap))
	if result > profile.DailyCapMax {
		return profile.DailyCapMax
	}
	return result
}

func (e *Engine) Jitter(profile *domain.ThrottleProfile, sentSincePause int) time.Duration {
	delay := rand.Intn(profile.MaxDelaySec-profile.MinDelaySec+1) + profile.MinDelaySec

	if profile.CoffeeBreakEveryN > 0 && sentSincePause >= profile.CoffeeBreakEveryN {
		extra := rand.Intn(profile.CoffeeBreakMaxSec-profile.CoffeeBreakMinSec+1) + profile.CoffeeBreakMinSec
		delay += extra
	}

	return time.Duration(delay) * time.Second
}

func (e *Engine) IsQuietHours(profile *domain.ThrottleProfile, recipientTZ string) bool {
	loc := time.UTC
	if profile.RespectTimezone && recipientTZ != "" {
		if parsed, err := time.LoadLocation(recipientTZ); err == nil {
			loc = parsed
		}
	}
	now := time.Now().In(loc)
	current := now.Hour()*60 + now.Minute()

	start := parseTimeMinutes(profile.QuietHoursStart)
	end := parseTimeMinutes(profile.QuietHoursEnd)

	if start > end {
		return current >= start || current < end
	}
	return current >= start && current < end
}

func (e *Engine) NextBusinessOpen(profile *domain.ThrottleProfile, recipientTZ string) time.Time {
	loc := time.UTC
	if profile.RespectTimezone && recipientTZ != "" {
		if parsed, err := time.LoadLocation(recipientTZ); err == nil {
			loc = parsed
		}
	}

	now := time.Now().In(loc)
	endMinutes := parseTimeMinutes(profile.QuietHoursEnd)
	next := time.Date(now.Year(), now.Month(), now.Day(), endMinutes/60, endMinutes%60, 0, 0, loc)
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next.UTC()
}

func bucketDayKey(instanceID uuid.UUID) string {
	day := time.Now().UTC().Format("2006-01-02")
	return fmt.Sprintf("cmp:bucket:%s:day:%s", instanceID, day)
}

func bucketHourKey(instanceID uuid.UUID) string {
	hour := time.Now().UTC().Format("2006-01-02T15")
	return fmt.Sprintf("cmp:bucket:%s:hour:%s", instanceID, hour)
}

func nextSlotKey(instanceID uuid.UUID) string {
	return fmt.Sprintf("cmp:next_slot:%s", instanceID)
}

func (e *Engine) BucketAllow(ctx context.Context, instanceID uuid.UUID, dailyCap, hourlyCap int) (bool, error) {
	dayKey := bucketDayKey(instanceID)
	dayCount, err := e.rdb.Get(ctx, dayKey).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if dayCount >= dailyCap {
		return false, nil
	}

	hourKey := bucketHourKey(instanceID)
	hourCount, err := e.rdb.Get(ctx, hourKey).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if hourCount >= hourlyCap {
		return false, nil
	}

	return true, nil
}

func (e *Engine) BucketConsume(ctx context.Context, instanceID uuid.UUID) error {
	pipe := e.rdb.Pipeline()

	dayKey := bucketDayKey(instanceID)
	pipe.Incr(ctx, dayKey)
	pipe.Expire(ctx, dayKey, 25*time.Hour)

	hourKey := bucketHourKey(instanceID)
	pipe.Incr(ctx, hourKey)
	pipe.Expire(ctx, hourKey, 2*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

func (e *Engine) GetNextSlot(ctx context.Context, instanceID uuid.UUID) (time.Time, error) {
	val, err := e.rdb.Get(ctx, nextSlotKey(instanceID)).Int64()
	if err == redis.Nil {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(val, 0).UTC(), nil
}

func (e *Engine) SetNextSlot(ctx context.Context, instanceID uuid.UUID, at time.Time) error {
	return e.rdb.Set(ctx, nextSlotKey(instanceID), at.Unix(), 24*time.Hour).Err()
}

func (e *Engine) SentSincePauseKey(instanceID uuid.UUID) string {
	return fmt.Sprintf("cmp:sent_since_pause:%s", instanceID)
}

func (e *Engine) GetSentSincePause(ctx context.Context, instanceID uuid.UUID) (int, error) {
	val, err := e.rdb.Get(ctx, e.SentSincePauseKey(instanceID)).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (e *Engine) IncrSentSincePause(ctx context.Context, instanceID uuid.UUID) (int, error) {
	key := e.SentSincePauseKey(instanceID)
	val, err := e.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	e.rdb.Expire(ctx, key, 24*time.Hour)
	return int(val), nil
}

func (e *Engine) ResetSentSincePause(ctx context.Context, instanceID uuid.UUID) error {
	return e.rdb.Del(ctx, e.SentSincePauseKey(instanceID)).Err()
}

func parseTimeMinutes(t string) int {
	var h, m int
	fmt.Sscanf(t, "%d:%d", &h, &m)
	return h*60 + m
}
