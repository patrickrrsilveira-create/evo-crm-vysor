package services

import (
	"context"
	"fmt"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/oauth"
	"gorm.io/gorm"
)

// OAuthService é o hub central de gerenciamento de fluxos OAuth
type OAuthService struct {
	providers map[string]oauth.Provider
	db        *gorm.DB
}

// NewOAuthService cria uma instância do hub OAuth
func NewOAuthService(db *gorm.DB) *OAuthService {
	return &OAuthService{
		providers: make(map[string]oauth.Provider),
		db:        db,
	}
}

// RegisterProvider adiciona uma estratégia (Google, HubSpot, etc) ao Hub
func (s *OAuthService) RegisterProvider(p oauth.Provider) {
	s.providers[p.Name()] = p
}

// GenerateAuthorizationURL cria a URL segura de login para o usuário
func (s *OAuthService) GenerateAuthorizationURL(agentID, providerName string) (string, error) {
	provider, exists := s.providers[providerName]
	if !exists {
		return "", fmt.Errorf("provedor '%s' não suportado", providerName)
	}

	// 1. Gerar PKCE Seguros
	verifier, challenge, err := oauth.GeneratePKCE()
	if err != nil {
		return "", err
	}

	// 2. O State também precisa ser seguro e forte (contra CSRF)
	stateBytes, _, _ := oauth.GeneratePKCE() // reaproveitando gerador para entropia
	state := stateBytes[:16]                 // string curta de 16 chars já serve para state

	// 3. Salvar sessão no banco de dados para recuperar no callback (substituindo mem-dict do Python)
	session := models.OAuthSession{
		State:     state,
		Provider:  providerName,
		AgentID:   agentID,
		Verifier:  verifier,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.db.Create(&session).Error; err != nil {
		return "", fmt.Errorf("falha ao persistir sessão OAuth: %w", err)
	}

	// 4. Delegar ao provedor específico a construção da URL
	authURL := provider.BuildAuthURL(state, challenge)

	return authURL, nil
}

// HandleCallback finaliza o login do usuário, trocando o código pelo Access Token
func (s *OAuthService) HandleCallback(ctx context.Context, state, code string) (*oauth.TokenResponse, error) {
	// 1. Resgatar a sessão no DB baseada no state
	var session models.OAuthSession
	if err := s.db.Where("state = ? AND expires_at > ?", state, time.Now()).First(&session).Error; err != nil {
		return nil, fmt.Errorf("sessão OAuth inválida, não encontrada ou expirada (CSRF prevent)")
	}

	// 2. Destruir o state para que não possa ser reusado (prevenção contra replay attacks)
	s.db.Delete(&session)

	// 3. Obter provedor
	provider, exists := s.providers[session.Provider]
	if !exists {
		return nil, fmt.Errorf("provedor inválido na sessão")
	}

	// 4. Delegar troca de token
	tokenResp, err := provider.ExchangeCode(ctx, code, session.Verifier)
	if err != nil {
		return nil, fmt.Errorf("falha ao trocar código pelo token: %w", err)
	}

	// 5. Aqui no futuro nós conectaremos com a gravação de `evo_core_agent_integrations`
	// Por enquanto, apenas retornamos o payload processado.

	return tokenResp, nil
}
