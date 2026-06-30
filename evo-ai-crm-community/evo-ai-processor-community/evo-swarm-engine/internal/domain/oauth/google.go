package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GoogleProvider implementa a interface Provider para todos os serviços Google.
// Diferencia a requisição de escopo baseado na ServiceVariant (e.g. google_calendar, google_sheets)
type GoogleProvider struct {
	ServiceVariant string
	ClientID       string
	ClientSecret   string
	RedirectURI    string
}

// NewGoogleProvider constrói a estratégia OAuth do Google, buscando credenciais do env
func NewGoogleProvider(variant string) *GoogleProvider {
	// A URL de callback é a rota centralizada unificada do Go Hub
	// Exemplo: http://localhost:8080/api/v1/oauth/callback/google_calendar
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	redirectURI := fmt.Sprintf("%s/api/v1/oauth/callback/%s", backendURL, variant)

	return &GoogleProvider{
		ServiceVariant: variant,
		ClientID:       os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURI:    redirectURI,
	}
}

// Name retorna o nome da variante, ex: "google_calendar", "google_drive"
func (g *GoogleProvider) Name() string {
	return g.ServiceVariant
}

// BuildAuthURL constrói a URL de login do Google requisitando os escopos corretos
func (g *GoogleProvider) BuildAuthURL(state, codeChallenge string) string {
	scopes := g.getScopesForVariant()

	authURL, _ := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	q := authURL.Query()
	q.Set("client_id", g.ClientID)
	q.Set("redirect_uri", g.RedirectURI)
	q.Set("response_type", "code")
	q.Set("scope", scopes)
	q.Set("state", state)
	q.Set("access_type", "offline")
	q.Set("prompt", "consent") // Garante que o Google sempre nos devolva um Refresh Token

	// Adicionar PKCE Security
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")

	authURL.RawQuery = q.Encode()
	return authURL.String()
}

// ExchangeCode troca o código de autorização pelo Access/Refresh Token na API do Google
func (g *GoogleProvider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("client_id", g.ClientID)
	data.Set("client_secret", g.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", g.RedirectURI)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição de token ao Google: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Google: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("google retornou erro %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON de tokens: %w", err)
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    result.TokenType,
	}, nil
}

// getScopesForVariant retorna as permissões (scopes) específicas para cada módulo do Workspace
func (g *GoogleProvider) getScopesForVariant() string {
	baseScopes := "openid email profile"

	switch g.ServiceVariant {
	case "google_calendar":
		return baseScopes + " https://www.googleapis.com/auth/calendar"
	case "google_sheets":
		return baseScopes + " https://www.googleapis.com/auth/spreadsheets"
	case "google_drive":
		// Usamos readonly para garantir segurança da IA, permitindo leitura de PDFs sem poder apagar arquivos de clientes.
		return baseScopes + " https://www.googleapis.com/auth/drive.readonly"
	default:
		return baseScopes
	}
}
