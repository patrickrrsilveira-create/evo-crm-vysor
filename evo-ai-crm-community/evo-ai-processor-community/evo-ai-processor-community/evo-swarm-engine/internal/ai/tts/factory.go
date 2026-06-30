package tts

import (
	"fmt"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)

// NewTTSProvider retorna o cliente correto baseado no provedor
func NewTTSProvider(providerName string, apiKey string) (models.TTSProvider, error) {
	if providerName == "elevenlabs" {
		return NewElevenLabsClient(apiKey), nil
	}

	// Adicionar fallback para outros provedores futuramente
	// if providerName == "google" { ... }

	return nil, fmt.Errorf("provedor TTS não suportado: %s", providerName)
}
