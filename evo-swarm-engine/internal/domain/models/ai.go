package models

import "context"

// LLMProvider define a interface base para qualquer LLM (OpenAI, Anthropic, Gemini, Groq)
type LLMProvider interface {
	Generate(ctx context.Context, req LLMRequest) (*LLMResponse, error)
	Name() string
}

// TTSProvider define a interface base para conversores de texto para fala
type TTSProvider interface {
	Synthesize(ctx context.Context, text string, options TTSOptions) ([]byte, error)
	Name() string
}

// LLMRequest encapsula a chamada para o provedor de IA
type LLMRequest struct {
	Model       string
	System      string
	Messages    []LLMMessage
	Temperature float32
	MaxTokens   int
	Tools       []LLMTool
}

type LLMMessage struct {
	Role       string        `json:"role"` // "system", "user", "assistant", "tool"
	Content    string        `json:"content"`
	Name       string        `json:"name,omitempty"`         // Opcional, usado para tool calls
	ToolCallID string        `json:"tool_call_id,omitempty"` // Opcional, para respostas de ferramentas
	ToolCalls  []LLMToolCall `json:"tool_calls,omitempty"`   // Quando a IA responde chamando ferramentas
}

// LLMResponse encapsula o retorno unificado dos LLMs
type LLMResponse struct {
	Content    string        `json:"content"`
	ToolCalls  []LLMToolCall `json:"tool_calls,omitempty"`
	TokenUsage TokenUsage    `json:"token_usage"`
}

// LLMToolCall representa uma intenção de execução de função da IA
type LLMToolCall struct {
	ID           string `json:"id"`
	FunctionName string `json:"function_name"`
	Arguments    string `json:"arguments"` // JSON arguments
}

// TokenUsage rastreia custos
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// LLMTool descreve uma ferramenta disponível para o modelo (Function Calling)
type LLMTool struct {
	Type     string // usually "function"
	Function LLMFunctionDef
}

type LLMFunctionDef struct {
	Name        string
	Description string
	Parameters  map[string]interface{} // JSON Schema representação dos parametros
}

// TTSOptions configurações específicas da voz
type TTSOptions struct {
	VoiceID string
	Model   string
	Speed   float32
}
