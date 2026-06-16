package tts

import (
	"fmt"
)

// Factory holds registered TTS providers
type Factory struct {
	providers map[string]Provider
}

// NewFactory creates a new TTS factory with default providers
func NewFactory() *Factory {
	f := &Factory{
		providers: make(map[string]Provider),
	}

	f.RegisterProvider(NewElevenLabsProvider())
	f.RegisterProvider(NewCartesiaProvider())
	f.RegisterProvider(NewOpenRouterProvider())

	return f
}

// RegisterProvider adds a new TTS provider
func (f *Factory) RegisterProvider(p Provider) {
	f.providers[p.Name()] = p
}

// GetProvider retrieves a provider by name
func (f *Factory) GetProvider(name string) (Provider, error) {
	p, exists := f.providers[name]
	if !exists {
		return nil, fmt.Errorf("provedor TTS não suportado: %s", name)
	}
	return p, nil
}
