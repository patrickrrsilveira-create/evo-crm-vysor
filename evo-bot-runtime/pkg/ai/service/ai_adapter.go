package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	brtErrors "github.com/EvolutionAPI/evo-bot-runtime/internal/errors"
	"github.com/EvolutionAPI/evo-bot-runtime/pkg/ai/model"
)

// maxResponseBytes caps the AI Processor response body to prevent OOM on oversized payloads.
const maxResponseBytes = 1 << 20 // 1 MiB

// AIAdapter calls the AI Processor via A2A protocol (JSON-RPC 2.0).
// Swap the backend by providing a different implementation at main.go wiring.
type AIAdapter interface {
	Call(ctx context.Context, req *model.A2ARequest) (*model.NormalizedResponse, error)
}

type aiAdapter struct {
	timeoutSecs int
	client      *http.Client
}

// NewAIAdapter constructs the adapter. Returns interface (GEAR R03).
// The AI Processor URL comes from each event's outgoing_url field.
func NewAIAdapter(timeoutSecs int) AIAdapter {
	return &aiAdapter{
		timeoutSecs: timeoutSecs,
		client:      &http.Client{},
	}
}

func (a *aiAdapter) Call(ctx context.Context, req *model.A2ARequest) (*model.NormalizedResponse, error) {
	start := time.Now()

	// Wrap with timeout — inner timeout, outer ctx for pipeline cancellation.
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(a.timeoutSecs)*time.Second)
	defer cancel()

	// Use the full outgoing_url provided by the CRM (already contains the agent ID)
	url := req.OutgoingURL

	// Build JSON-RPC 2.0 envelope
	rpcReq := model.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      fmt.Sprintf("%d:%d", req.ContactID, req.ConversationID),
		Method:  "message/send",
		Params: model.JSONRPCParams{
			ContextID: fmt.Sprintf("%d", req.ConversationID),
			UserID:    fmt.Sprintf("%d", req.ContactID),
			Message: model.JSONRPCMessage{
				Role: "user",
				Parts: []model.JSONRPCPart{
					{Type: "text", Text: req.Message},
				},
			},
			Metadata: nonNilMetadata(req.Metadata),
		},
	}

	body, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("pipeline.ai.marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(timeoutCtx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("pipeline.ai.new_request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", req.ApiKey)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		if ctx.Err() != nil {
			return nil, brtErrors.ErrPipelineCancelled
		}
		if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
			slog.Warn("pipeline.ai.http.timeout",
				"contact_id", req.ContactID,
				"conversation_id", req.ConversationID,
				"timeout_secs", a.timeoutSecs,
			)
			return nil, brtErrors.ErrAITimeout
		}
		return nil, fmt.Errorf("pipeline.ai.http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pipeline.ai.status: unexpected %d from AI Processor", resp.StatusCode)
	}

	var a2aResp model.A2AResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes)).Decode(&a2aResp); err != nil {
		return nil, fmt.Errorf("pipeline.ai.decode: %w", err)
	}

	content := extractResponseText(&a2aResp)

	// Download audio payload from artifacts if available
	audioPayload, err := downloadAudioFromArtifacts(ctx, a.client, a.timeoutSecs, &a2aResp)
	if err != nil {
		slog.Error("pipeline.ai.download_audio.failed", "error", err)
	}

	slog.Info("pipeline.ai.http.completed",
		"contact_id", req.ContactID,
		"conversation_id", req.ConversationID,
		"duration_ms", time.Since(start).Milliseconds(),
		"generated_audio", audioPayload != nil,
	)

	return &model.NormalizedResponse{
		Content: content,
		Audio:   audioPayload,
	}, nil
}

// extractResponseText extracts the text content from the A2A JSON-RPC response.
// Tries result.artifacts[0].parts[0].text first, then result.message.parts[0].text.
func extractResponseText(resp *model.A2AResponse) string {
	if resp.Result == nil {
		return ""
	}
	// Try artifacts first (primary response format)
	if len(resp.Result.Artifacts) > 0 {
		for _, artifact := range resp.Result.Artifacts {
			for _, part := range artifact.Parts {
				if part.Text != "" {
					return part.Text
				}
			}
		}
	}
	// Fallback to message format
	if resp.Result.Message != nil {
		for _, part := range resp.Result.Message.Parts {
			if part.Text != "" {
				return part.Text
			}
		}
	}
	return ""
}

// nonNilMetadata ensures metadata is never nil (avoids "null" in JSON).
func nonNilMetadata(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}

// downloadAudioFromArtifacts checks the response for an audio URL in artifacts
// and downloads the audio payload.
func downloadAudioFromArtifacts(ctx context.Context, client *http.Client, timeoutSecs int, resp *model.A2AResponse) ([]byte, error) {
	if resp.Result == nil {
		return nil, nil
	}

	var audioURL string
	// Find the first audio artifact URL
	for _, artifact := range resp.Result.Artifacts {
		for _, part := range artifact.Parts {
			if part.Type == "file" && part.URL != "" && (part.MimeType == "audio/ogg" || part.MimeType == "audio/mpeg" || part.MimeType == "audio/opus") {
				audioURL = part.URL
				break
			}
		}
		if audioURL != "" {
			break
		}
	}

	if audioURL == "" {
		return nil, nil
	}

	slog.Info("pipeline.ai.download_audio.started", "url", audioURL)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, audioURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new audio req: %w", err)
	}

	downloadResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do audio req: %w", err)
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audio download bad status: %d", downloadResp.StatusCode)
	}

	audioData, err := io.ReadAll(downloadResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read audio resp: %w", err)
	}

	slog.Info("pipeline.ai.download_audio.completed", "bytes", len(audioData))
	return audioData, nil
}
