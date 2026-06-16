package tts

import (
	"context"
)

// Request encapsula os parâmetros de geração de áudio
type Request struct {
	Text     string // Texto a ser falado
	APIKey   string // Token de acesso
	VoiceID  string // ID da voz selecionada
	Model    string // Modelo de IA (ex: kokoro-82m, eleven_multilingual_v2)
	Language string // Idioma (pt, en, etc)
}

// Provider define o contrato para os motores de síntese de voz
type Provider interface {
	Name() string
	GenerateSpeech(ctx context.Context, req Request) ([]byte, error)
}
