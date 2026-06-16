package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CartesiaProvider implements TTS for Cartesia
type CartesiaProvider struct{}

func NewCartesiaProvider() *CartesiaProvider {
	return &CartesiaProvider{}
}

func (p *CartesiaProvider) Name() string {
	return "cartesia"
}

func (p *CartesiaProvider) GenerateSpeech(ctx context.Context, req Request) ([]byte, error) {
	if req.APIKey == "" || req.VoiceID == "" {
		return nil, fmt.Errorf("cartesia requires APIKey and VoiceID")
	}

	url := "https://api.cartesia.ai/tts/bytes"

	payload := map[string]interface{}{
		"model_id":   "sonic-english",
		"transcript": req.Text,
		"voice": map[string]interface{}{
			"mode": "id",
			"id":   req.VoiceID,
		},
		"output_format": map[string]interface{}{
			"container":   "raw",
			"encoding":    "pcm_s16le",
			"sample_rate": 24000,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("X-API-Key", req.APIKey)
	httpReq.Header.Set("Cartesia-Version", "2024-06-10")
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cartesia request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorDetail, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cartesia API error: %d - %s", resp.StatusCode, string(errorDetail))
	}

	return io.ReadAll(resp.Body)
}
