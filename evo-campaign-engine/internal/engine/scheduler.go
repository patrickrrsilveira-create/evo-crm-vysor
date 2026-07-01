package engine

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/repository"
	"log"
	"time"
)

type Scheduler struct {
	jobRepo      *repository.SendJobRepo
	campaignRepo *repository.CampaignRepo
	tick         time.Duration
	batchSize    int
	jobCh        chan domain.SendJob
}

func NewScheduler(
	jobRepo *repository.SendJobRepo,
	campaignRepo *repository.CampaignRepo,
	tick time.Duration,
	batchSize int,
	jobCh chan domain.SendJob,
) *Scheduler {
	return &Scheduler{
		jobRepo:      jobRepo,
		campaignRepo: campaignRepo,
		tick:         tick,
		batchSize:    batchSize,
		jobCh:        jobCh,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	log.Printf("scheduler started (tick=%s batch=%d)", s.tick, s.batchSize)
	ticker := time.NewTicker(s.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler stopped")
			return
		case <-ticker.C:
			s.promote(ctx)
		}
	}
}

func (s *Scheduler) promote(ctx context.Context) {
	jobs, err := s.jobRepo.FetchDueJobs(ctx, s.batchSize)
	if err != nil {
		log.Printf("scheduler: fetch due jobs error: %v", err)
		return
	}
	if len(jobs) == 0 {
		return
	}

	for _, job := range jobs {
		if err := s.jobRepo.UpdateState(ctx, job.ID, domain.JobScheduled, nil); err != nil {
			log.Printf("scheduler: update job %s state error: %v", job.ID, err)
			continue
		}

		select {
		case s.jobCh <- job:
		case <-ctx.Done():
			return
		}
	}

	log.Printf("scheduler: promoted %d jobs", len(jobs))
}
