package engine

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/driver"
	"evo-campaign-engine/internal/repository"
	"evo-campaign-engine/internal/spintax"
	"evo-campaign-engine/internal/throttle"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	id             int
	jobRepo        *repository.SendJobRepo
	campaignRepo   *repository.CampaignRepo
	senderRepo     *repository.SenderRepo
	suppressRepo   *repository.SuppressionRepo
	throttleRepo   *repository.ThrottleRepo
	throttleEngine *throttle.Engine
	rotator        *Rotator
	drivers        map[string]driver.ChannelDriver
	maxRetries     int
	jobCh          <-chan domain.SendJob
}

type WorkerDeps struct {
	JobRepo        *repository.SendJobRepo
	CampaignRepo   *repository.CampaignRepo
	SenderRepo     *repository.SenderRepo
	SuppressRepo   *repository.SuppressionRepo
	ThrottleRepo   *repository.ThrottleRepo
	ThrottleEngine *throttle.Engine
	Rotator        *Rotator
	Drivers        map[string]driver.ChannelDriver
	MaxRetries     int
}

func NewWorker(id int, deps WorkerDeps, jobCh <-chan domain.SendJob) *Worker {
	return &Worker{
		id:             id,
		jobRepo:        deps.JobRepo,
		campaignRepo:   deps.CampaignRepo,
		senderRepo:     deps.SenderRepo,
		suppressRepo:   deps.SuppressRepo,
		throttleRepo:   deps.ThrottleRepo,
		throttleEngine: deps.ThrottleEngine,
		rotator:        deps.Rotator,
		drivers:        deps.Drivers,
		maxRetries:     deps.MaxRetries,
		jobCh:          jobCh,
	}
}

func (w *Worker) Run(ctx context.Context) {
	log.Printf("worker-%d started", w.id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker-%d stopped", w.id)
			return
		case job := <-w.jobCh:
			w.process(ctx, job)
		}
	}
}

func (w *Worker) process(ctx context.Context, job domain.SendJob) {
	campaign, err := w.campaignRepo.GetByID(ctx, job.CampaignID)
	if err != nil {
		log.Printf("worker-%d: campaign %s not found: %v", w.id, job.CampaignID, err)
		w.fail(ctx, job, "campaign not found")
		return
	}

	if campaign.Status != domain.CampaignRunning {
		w.skip(ctx, job)
		return
	}

	suppressed, _ := w.suppressRepo.IsSuppressed(ctx, campaign.AccountID, job.ContactID, "")
	if suppressed {
		w.jobRepo.UpdateState(ctx, job.ID, domain.JobSkipped, map[string]interface{}{"last_error": "suppressed"})
		return
	}

	channels, err := w.campaignRepo.GetChannels(ctx, job.CampaignID)
	if err != nil || len(channels) == 0 {
		w.reschedule(ctx, job, 60*time.Second)
		return
	}

	inboxIDs := make([]int, len(channels))
	for i, ch := range channels {
		inboxIDs[i] = ch.InboxID
	}

	sender, err := w.rotator.Pick(ctx, inboxIDs)
	if err != nil || sender == nil {
		w.reschedule(ctx, job, 60*time.Second)
		return
	}

	var profile *domain.ThrottleProfile
	if campaign.ThrottleProfileID != nil {
		profile, _ = w.throttleRepo.GetByID(ctx, *campaign.ThrottleProfileID)
	}
	if profile == nil {
		profile = defaultProfile()
	}

	if w.throttleEngine.IsQuietHours(profile, "") {
		nextOpen := w.throttleEngine.NextBusinessOpen(profile, "")
		w.reschedule(ctx, job, time.Until(nextOpen))
		return
	}

	nextSlot, _ := w.throttleEngine.GetNextSlot(ctx, sender.ID)
	if !nextSlot.IsZero() && time.Now().UTC().Before(nextSlot) {
		w.reschedule(ctx, job, time.Until(nextSlot))
		return
	}

	warmupDay := sender.WarmupDay
	dailyCap := w.throttleEngine.WarmupCap(profile, warmupDay)
	allowed, _ := w.throttleEngine.BucketAllow(ctx, sender.ID, dailyCap, profile.HourlyCap)
	if !allowed {
		w.reschedule(ctx, job, 10*time.Minute)
		return
	}

	w.jobRepo.AssignSender(ctx, job.ID, sender.ID)
	w.jobRepo.UpdateState(ctx, job.ID, domain.JobSending, nil)

	variant := w.pickVariant(ctx, campaign)
	body := ""
	mediaURL := ""
	mediaType := ""
	if variant != nil {
		body = spintax.Process(variant.Body)
		mediaURL = variant.MediaURL
		mediaType = variant.MediaType
	}

	drv := w.resolveDriver(sender.ChannelType)
	if drv == nil {
		w.fail(ctx, job, "no driver for channel: "+sender.ChannelType)
		return
	}

	content := driver.Content{
		Text:      body,
		MediaURL:  mediaURL,
		MediaType: mediaType,
		Instance:  sender.Identifier,
	}

	result := drv.Send(ctx, sender.Identifier, job.Recipient, content)

	w.saveEvent(ctx, job.ID, result)

	if result.OK {
		now := time.Now().UTC()
		w.jobRepo.UpdateState(ctx, job.ID, domain.JobSent, map[string]interface{}{"sent_at": now})
		w.campaignRepo.IncrementCounter(ctx, job.CampaignID, "sent_count", 1)
		w.throttleEngine.BucketConsume(ctx, sender.ID)
		w.senderRepo.IncrementSentToday(ctx, sender.ID)

		sentSincePause, _ := w.throttleEngine.IncrSentSincePause(ctx, sender.ID)
		jitter := w.throttleEngine.Jitter(profile, sentSincePause)

		if profile.CoffeeBreakEveryN > 0 && sentSincePause >= profile.CoffeeBreakEveryN {
			w.throttleEngine.ResetSentSincePause(ctx, sender.ID)
		}

		w.throttleEngine.SetNextSlot(ctx, sender.ID, now.Add(jitter))
	} else {
		job.Attempts++
		if job.Attempts < w.maxRetries {
			backoff := time.Duration(job.Attempts*30) * time.Second
			w.jobRepo.UpdateState(ctx, job.ID, domain.JobScheduled, map[string]interface{}{
				"attempts":     job.Attempts,
				"last_error":   result.Error,
				"scheduled_at": time.Now().UTC().Add(backoff),
			})
		} else {
			w.fail(ctx, job, result.Error)
			w.campaignRepo.IncrementCounter(ctx, job.CampaignID, "failed_count", 1)
		}
	}
}

func (w *Worker) pickVariant(ctx context.Context, campaign *domain.Campaign) *domain.MessageVariant {
	if len(campaign.Variants) == 0 {
		return nil
	}

	totalWeight := 0
	for _, v := range campaign.Variants {
		totalWeight += v.Weight
	}
	if totalWeight == 0 {
		return &campaign.Variants[0]
	}

	r := rand.Intn(totalWeight)
	cumulative := 0
	for i := range campaign.Variants {
		cumulative += campaign.Variants[i].Weight
		if r < cumulative {
			return &campaign.Variants[i]
		}
	}
	return &campaign.Variants[0]
}

func (w *Worker) resolveDriver(channelType string) driver.ChannelDriver {
	if drv, ok := w.drivers[channelType]; ok {
		return drv
	}
	if drv, ok := w.drivers["whatsapp"]; ok {
		return drv
	}
	return nil
}

func (w *Worker) fail(ctx context.Context, job domain.SendJob, reason string) {
	w.jobRepo.UpdateState(ctx, job.ID, domain.JobFailed, map[string]interface{}{"last_error": reason})
}

func (w *Worker) skip(ctx context.Context, job domain.SendJob) {
	w.jobRepo.UpdateState(ctx, job.ID, domain.JobSkipped, nil)
}

func (w *Worker) reschedule(ctx context.Context, job domain.SendJob, delay time.Duration) {
	w.jobRepo.Reschedule(ctx, job.ID, time.Now().UTC().Add(delay))
}

func (w *Worker) saveEvent(ctx context.Context, jobID uuid.UUID, result driver.SendResult) {
	status := "sent"
	if !result.OK {
		status = "failed"
	}
	w.jobRepo.SaveDeliveryEvent(ctx, &domain.DeliveryEvent{
		SendJobID:      jobID,
		Status:         status,
		ProviderStatus: result.ProviderStatus,
		Ts:             time.Now().UTC(),
	})
}

func defaultProfile() *domain.ThrottleProfile {
	return &domain.ThrottleProfile{
		DailyCapStart:     20,
		DailyCapMax:       500,
		WarmupMultiplier:  1.3,
		WarmupStepDays:    1,
		HourlyCap:         60,
		MinDelaySec:       45,
		MaxDelaySec:       120,
		CoffeeBreakEveryN: 25,
		CoffeeBreakMinSec: 120,
		CoffeeBreakMaxSec: 300,
		QuietHoursStart:   "20:00",
		QuietHoursEnd:     "08:00",
		RespectTimezone:   true,
	}
}
