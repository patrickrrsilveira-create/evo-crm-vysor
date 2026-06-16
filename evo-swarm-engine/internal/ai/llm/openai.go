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

type OpenAIClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewOpenAIClient(apiKey string, baseURL string) *OpenAIClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	}
	return &OpenAIClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *OpenAIClient) Name() string {
	return "openai"
}

// openAIRequest e openAIResponse são as estruturas locais da API deles
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float32         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Tools       []openAITool    `json:"tools,omitempty"`
}

type openAIMessage struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	Name       string `json:"name,omitempty"`
	ToolCallID string `json:"tool_call_id,omitempty"`
}

type openAITool struct {
	Type     string         `json:"type"`
	Function openAIFunction `json:"function"`
}

type openAIFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *OpenAIClient) Generate(ctx context.Context, req models.LLMRequest) (*models.LLMResponse, error) {
	modelName := req.Model
	if strings.HasPrefix(modelName, "openrouter/") {
		modelName = strings.TrimPrefix(modelName, "openrouter/")
	}

	openReq := openAIRequest{
		Model:       modelName,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	if req.System != "" {
		openReq.Messages = append(openReq.Messages, openAIMessage{Role: "system", Content: req.System})
	}

	for _, msg := range req.Messages {
		openReq.Messages = append(openReq.Messages, openAIMessage{
			Role:       msg.Role,
			Content:    msg.Content,
			Name:       msg.Name,
			ToolCallID: msg.ToolCallID,
		})
	}

	for _, t := range req.Tools {
		openReq.Tools = append(openReq.Tools, openAITool{
			Type: "function",
			Function: openAIFunction{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		})
	}

	// 2. Serializar JSON e disparar HTTP POST
	payloadBytes, err := json.Marshal(openReq)
	if err != nil {
		return nil, fmt.Errorf("erro serializando payload OpenAI: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro chamando API OpenAI: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro OpenAI API (status %d): %s", resp.StatusCode, string(body))
	}

	var openResp openAIResponse
	if err := json.Unmarshal(body, &openResp); err != nil {
		return nil, fmt.Errorf("erro unmarshalling resposta OpenAI: %v", err)
	}

	if len(openResp.Choices) == 0 {
		return nil, fmt.Errorf("openai retornou 0 choices")
	}

	choice := openResp.Choices[0].Message

	// 3. Converter retorno
	out := &models.LLMResponse{
		Content: choice.Content,
		TokenUsage: models.TokenUsage{
			PromptTokens:     openResp.Usage.PromptTokens,
			CompletionTokens: openResp.Usage.CompletionTokens,
			TotalTokens:      openResp.Usage.TotalTokens,
		},
	}

	for _, tc := range choice.ToolCalls {
		out.ToolCalls = append(out.ToolCalls, models.LLMToolCall{
			ID:           tc.ID,
			FunctionName: tc.Function.Name,
			Arguments:    tc.Function.Arguments,
		})
	}

	return out, nil
}
