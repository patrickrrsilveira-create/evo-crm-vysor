package driver

import "context"

type SendResult struct {
	OK             bool
	ProviderStatus string
	Error          string
}

type Content struct {
	Text      string
	MediaURL  string
	MediaType string
	Instance  string
}

type ChannelDriver interface {
	Send(ctx context.Context, instance, recipient string, content Content) SendResult
}
