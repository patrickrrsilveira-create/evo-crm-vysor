package engine

import (
	"context"
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/repository"
	"sync/atomic"
)

type Rotator struct {
	senderRepo *repository.SenderRepo
	index      atomic.Uint64
}

func NewRotator(senderRepo *repository.SenderRepo) *Rotator {
	return &Rotator{senderRepo: senderRepo}
}

func (r *Rotator) Pick(ctx context.Context, inboxIDs []int) (*domain.SenderInstance, error) {
	senders, err := r.senderRepo.GetByInboxIDs(ctx, inboxIDs)
	if err != nil {
		return nil, err
	}
	if len(senders) == 0 {
		return nil, nil
	}

	idx := r.index.Add(1)
	picked := senders[int(idx)%len(senders)]
	return &picked, nil
}
