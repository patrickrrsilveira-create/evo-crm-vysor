package tts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// OpenRouterProvider implements TTS for OpenRouter (e.g. Kokoro)
type OpenRouterProvider struct{}

func NewOpenRouterProvider() *OpenRouterProvider {
	return &OpenRouterProvider{}
}

func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

func (p *OpenRouterProvider) GenerateSpeech(ctx context.Context, req Request) ([]byte, error) {
	if req.APIKey == "" {
		return nil, fmt.Errorf("openrouter requires APIKey")
	}

	model := req.Model
	if model == "" {
		model = "hexgrad/kokoro-82m"
	}

	voice := req.VoiceID
	if voice == "" {
		voice = "af_heart"
	}

	url := "https://openrouter.ai/api/v1/audio/speech"

	language := req.Language
	if language == "" {
		language = "pt"
	}

	payload := map[string]interface{}{
		"input":    req.Text,
		"model":    model,
		"voice":    voice,
		"language": language,
		"lang":     language,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openrouter request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorDetail, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openrouter API error: %d - %s", resp.StatusCode, string(errorDetail))
	}

	contentType := resp.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read openrouter response: %w", err)
	}

	// 1. JSON payload
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("failed to parse json response: %w", err)
		}
		if errMsg, ok := result["error"]; ok {
			return nil, fmt.Errorf("openrouter returned error: %v", errMsg)
		}

		if audio, ok := result["audio"].(string); ok {
			return base64.StdEncoding.DecodeString(audio)
		}
		if dataStr, ok := result["data"].(string); ok {
			return base64.StdEncoding.DecodeString(dataStr)
		}
		if dataArr, ok := result["data"].([]interface{}); ok && len(dataArr) > 0 {
			if firstItem, ok := dataArr[0].(map[string]interface{}); ok {
				if b64Str, ok := firstItem["b64_json"].(string); ok {
					return base64.StdEncoding.DecodeString(b64Str)
				}
			}
		}
		return nil, fmt.Errorf("openrouter returned json but no audio payload found")
	}

	// 2. PCM Stream (Wrap in WAV)
	if strings.Contains(strings.ToLower(contentType), "audio/pcm") {
		rate := 24000
		channels := 1

		rateRe := regexp.MustCompile(`rate=(\d+)`)
		if matches := rateRe.FindStringSubmatch(strings.ToLower(contentType)); len(matches) > 1 {
			if r, err := strconv.Atoi(matches[1]); err == nil {
				rate = r
			}
		}

		chRe := regexp.MustCompile(`channels=(\d+)`)
		if matches := chRe.FindStringSubmatch(strings.ToLower(contentType)); len(matches) > 1 {
			if c, err := strconv.Atoi(matches[1]); err == nil {
				channels = c
			}
		}

		return wrapPCMInWAV(bodyBytes, rate, channels)
	}

	return bodyBytes, nil
}

// wrapPCMInWAV converts raw PCM bytes to a valid WAV file in memory
func wrapPCMInWAV(pcm []byte, sampleRate int, numChannels int) ([]byte, error) {
	bitsPerSample := 16
	byteRate := sampleRate * numChannels * (bitsPerSample / 8)
	blockAlign := numChannels * (bitsPerSample / 8)
	dataSize := len(pcm)
	chunkSize := 36 + dataSize

	buf := new(bytes.Buffer)

	// RIFF header
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, uint32(chunkSize))
	buf.WriteString("WAVE")

	// fmt subchunk
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16)) // Subchunk1Size for PCM
	binary.Write(buf, binary.LittleEndian, uint16(1))  // AudioFormat (1 = PCM)
	binary.Write(buf, binary.LittleEndian, uint16(numChannels))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(buf, binary.LittleEndian, uint32(byteRate))
	binary.Write(buf, binary.LittleEndian, uint16(blockAlign))
	binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample))

	// data subchunk
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))

	// Append PCM data
	buf.Write(pcm)

	return buf.Bytes(), nil
}
