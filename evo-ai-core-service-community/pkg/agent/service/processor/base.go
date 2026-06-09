package processor

import (
	"context"

	"evo-ai-core-service/pkg/agent/model"
)

type BaseProcessor interface {
	Create(ctx context.Context, agent *model.Agent) error
	Update(ctx context.Context, current, request *model.Agent) error
}
