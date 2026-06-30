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

type GeminiClient struct {
	APIKey string
	Client *http.Client
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		APIKey: apiKey,
		Client: &http.Client{},
	}
}

func (c *GeminiClient) Name() string {
	return "gemini"
}

// geminiRequest e geminiResponse
type geminiRequest struct {
	Contents          []geminiContent `json:"contents"`
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"`
	Tools             []geminiTool    `json:"tools,omitempty"`
	GenerationConfig  struct {
		Temperature float32 `json:"temperature,omitempty"`
		MaxTokens   int     `json:"maxOutputTokens,omitempty"`
	} `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text         string              `json:"text,omitempty"`
	FunctionCall *geminiFunctionCall `json:"functionCall,omitempty"`
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDeclaration `json:"functionDeclarations"`
}

type geminiFunctionDeclaration struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

func (c *GeminiClient) Generate(ctx context.Context, req models.LLMRequest) (*models.LLMResponse, error) {
	gemReq := geminiRequest{}
	gemReq.GenerationConfig.Temperature = req.Temperature
	gemReq.GenerationConfig.MaxTokens = req.MaxTokens

	if req.System != "" {
		gemReq.SystemInstruction = &geminiContent{
			Role:  "system",
			Parts: []geminiPart{{Text: req.System}},
		}
	}

	for _, msg := range req.Messages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		} else if msg.Role == "tool" {
			role = "function"
		}

		gemReq.Contents = append(gemReq.Contents, geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: msg.Content}},
		})
	}

	if len(req.Tools) > 0 {
		tool := geminiTool{}
		for _, t := range req.Tools {
			tool.FunctionDeclarations = append(tool.FunctionDeclarations, geminiFunctionDeclaration{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			})
		}
		gemReq.Tools = append(gemReq.Tools, tool)
	}

	payloadBytes, err := json.Marshal(gemReq)
	if err != nil {
		return nil, fmt.Errorf("erro serializando payload Gemini: %v", err)
	}

	// Use v1beta for tools support
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", req.Model, c.APIKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro chamando API Gemini: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro Gemini API (status %d): %s", resp.StatusCode, string(body))
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(body, &gemResp); err != nil {
		return nil, fmt.Errorf("erro unmarshalling resposta Gemini: %v", err)
	}

	if len(gemResp.Candidates) == 0 {
		return nil, fmt.Errorf("gemini retornou 0 candidates")
	}

	out := &models.LLMResponse{
		TokenUsage: models.TokenUsage{
			PromptTokens:     gemResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: gemResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      gemResp.UsageMetadata.TotalTokenCount,
		},
	}

	for _, part := range gemResp.Candidates[0].Content.Parts {
		if part.Text != "" {
			out.Content += part.Text
		}
		if part.FunctionCall != nil {
			argsJson, _ := json.Marshal(part.FunctionCall.Args)
			out.ToolCalls = append(out.ToolCalls, models.LLMToolCall{
				ID:           part.FunctionCall.Name, // Gemini uses name as ID essentially
				FunctionName: part.FunctionCall.Name,
				Arguments:    string(argsJson),
			})
		}
	}

	return out, nil
}
