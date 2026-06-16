package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

type ElevenLabsClient struct {
	APIKey string
	Client *http.Client
}

func NewElevenLabsClient(apiKey string) *ElevenLabsClient {
	return &ElevenLabsClient{
		APIKey: apiKey,
		Client: &http.Client{},
	}
}

func (c *ElevenLabsClient) Name() string {
	return "elevenlabs"
}

type elevenLabsRequest struct {
	Text    string `json:"text"`
	ModelID string `json:"model_id,omitempty"`
}

func (c *ElevenLabsClient) Synthesize(ctx context.Context, text string, options models.TTSOptions) ([]byte, error) {
	if options.VoiceID == "" {
		return nil, fmt.Errorf("VoiceID é obrigatório para ElevenLabs")
	}

	modelID := options.Model
	if modelID == "" {
		modelID = "eleven_multilingual_v2"
	}

	reqBody := elevenLabsRequest{
		Text:    text,
		ModelID: modelID,
	}

	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro serializando payload ElevenLabs: %v", err)
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", options.VoiceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("xi-api-key", c.APIKey)
	httpReq.Header.Set("Accept", "audio/mpeg")

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro chamando API ElevenLabs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ElevenLabs API (status %d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
