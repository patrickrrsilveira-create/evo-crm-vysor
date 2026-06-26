package oauth

import "context"

// TokenResponse representa o payload de resposta padrão do provedor OAuth 2.0
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ConfigProvider define a configuração base de um provedor OAuth
type ConfigProvider struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURI  string
	Scopes       []string
}

// Provider define o contrato (Strategy Pattern) para qualquer integração OAuth (Google, Notion, etc)
type Provider interface {
	// Name retorna o identificador do provedor (ex: "google_calendar")
	Name() string

	// BuildAuthURL constrói a URL de login do usuário, injetando o State e PKCE
	BuildAuthURL(state, challenge string) string

	// ExchangeCode troca o código de autorização pelo AccessToken usando o verifier
	ExchangeCode(ctx context.Context, code, verifier string) (*TokenResponse, error)
}
