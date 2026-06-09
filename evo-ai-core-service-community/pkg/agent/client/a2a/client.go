package a2a

import (
	"context"
	"evo-ai-core-service/internal/httpclient"
	"fmt"
	"net/http"
)

type Client interface {
	FetchAgentCard(ctx context.Context, cardURL string) (map[string]interface{}, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) FetchAgentCard(ctx context.Context, cardURL string) (map[string]interface{}, error) {
	response, err := httpclient.DoGetJSON[map[string]interface{}](ctx, cardURL, nil, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent card: %v", err)
	}

	return *response, nil
}
