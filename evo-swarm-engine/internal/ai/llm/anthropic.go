package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

type AnthropicClient struct {
	APIKey string
	Client *http.Client
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		APIKey: apiKey,
		Client: &http.Client{},
	}
}

func (c *AnthropicClient) Name() string {
	return "anthropic"
}

// anthropicRequest e anthropicResponse
type anthropicRequest struct {
	Model       string             `json:"model"`
	System      string             `json:"system,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float32            `json:"temperature,omitempty"`
	Tools       []anthropicTool    `json:"tools,omitempty"`
}

type anthropicMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // String or array of content blocks
}

type anthropicTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type  string                 `json:"type"`
		Text  string                 `json:"text,omitempty"`
		ID    string                 `json:"id,omitempty"`
		Name  string                 `json:"name,omitempty"`
		Input map[string]interface{} `json:"input,omitempty"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func (c *AnthropicClient) Generate(ctx context.Context, req models.LLMRequest) (*models.LLMResponse, error) {
	antReq := anthropicRequest{
		Model:       req.Model,
		System:      req.System,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	if antReq.MaxTokens == 0 {
		antReq.MaxTokens = 4096 // Anthropic requires max_tokens
	}

	for _, msg := range req.Messages {
		antReq.Messages = append(antReq.Messages, anthropicMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	for _, t := range req.Tools {
		antReq.Tools = append(antReq.Tools, anthropicTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			InputSchema: t.Function.Parameters,
		})
	}

	payloadBytes, err := json.Marshal(antReq)
	if err != nil {
		return nil, fmt.Errorf("erro serializando payload Anthropic: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro chamando API Anthropic: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro Anthropic API (status %d): %s", resp.StatusCode, string(body))
	}

	var antResp anthropicResponse
	if err := json.Unmarshal(body, &antResp); err != nil {
		return nil, fmt.Errorf("erro unmarshalling resposta Anthropic: %v", err)
	}

	out := &models.LLMResponse{
		TokenUsage: models.TokenUsage{
			PromptTokens:     antResp.Usage.InputTokens,
			CompletionTokens: antResp.Usage.OutputTokens,
			TotalTokens:      antResp.Usage.InputTokens + antResp.Usage.OutputTokens,
		},
	}

	for _, block := range antResp.Content {
		if block.Type == "text" {
			out.Content += block.Text
		} else if block.Type == "tool_use" {
			argsJson, _ := json.Marshal(block.Input)
			out.ToolCalls = append(out.ToolCalls, models.LLMToolCall{
				ID:           block.ID,
				FunctionName: block.Name,
				Arguments:    string(argsJson),
			})
		}
	}

	return out, nil
}
