package engine

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"log"
	"sync"
)

type Pool struct {
	scheduler *Scheduler
	workers   []*Worker
	jobCh     chan domain.SendJob
}

func NewPool(scheduler *Scheduler, deps WorkerDeps, workerCount int) *Pool {
	jobCh := make(chan domain.SendJob, workerCount*10)

	scheduler.jobCh = jobCh

	workers := make([]*Worker, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = NewWorker(i, deps, jobCh)
	}

	return &Pool{
		scheduler: scheduler,
		workers:   workers,
		jobCh:     jobCh,
	}
}

func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.scheduler.Run(ctx)
	}()

	for _, w := range p.workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.Run(ctx)
		}(w)
	}

	log.Printf("engine pool started: 1 scheduler + %d workers", len(p.workers))
	wg.Wait()
	log.Println("engine pool stopped")
}
