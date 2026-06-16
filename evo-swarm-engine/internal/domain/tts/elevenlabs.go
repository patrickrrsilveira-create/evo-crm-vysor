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

// ElevenLabsProvider implements TTS for ElevenLabs
type ElevenLabsProvider struct{}

func NewElevenLabsProvider() *ElevenLabsProvider {
	return &ElevenLabsProvider{}
}

func (p *ElevenLabsProvider) Name() string {
	return "elevenlabs"
}

func (p *ElevenLabsProvider) GenerateSpeech(ctx context.Context, req Request) ([]byte, error) {
	if req.APIKey == "" || req.VoiceID == "" {
		return nil, fmt.Errorf("elevenlabs requires APIKey and VoiceID")
	}

	modelID := req.Model
	if modelID == "" {
		modelID = "eleven_multilingual_v2"
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", req.VoiceID)

	payload := map[string]interface{}{
		"text":     req.Text,
		"model_id": modelID,
		"voice_settings": map[string]interface{}{
			"stability":         0.8,
			"similarity_boost":  0.5,
			"style":             0.0,
			"use_speaker_boost": true,
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

	httpReq.Header.Set("xi-api-key", req.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("elevenlabs request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorDetail, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("elevenlabs API error: %d - %s", resp.StatusCode, string(errorDetail))
	}

	return io.ReadAll(resp.Body)
}
