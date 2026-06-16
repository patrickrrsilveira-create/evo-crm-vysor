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

// StripeProvider implementa a integração OAuth para o Stripe Connect
type StripeProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func NewStripeProvider() *StripeProvider {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}
	redirectURI := fmt.Sprintf("%s/api/v1/oauth/callback/stripe", backendURL)

	return &StripeProvider{
		ClientID:     os.Getenv("STRIPE_CLIENT_ID"),
		ClientSecret: os.Getenv("STRIPE_SECRET_KEY"), // Stripe usa a secret_key como client_secret
		RedirectURI:  redirectURI,
	}
}

func (s *StripeProvider) Name() string {
	return "stripe"
}

// BuildAuthURL constrói a URL de login do Stripe Connect
func (s *StripeProvider) BuildAuthURL(state, codeChallenge string) string {
	authURL, _ := url.Parse("https://connect.stripe.com/oauth/authorize")
	q := authURL.Query()
	q.Set("client_id", s.ClientID)
	q.Set("response_type", "code")
	q.Set("scope", "read_write")
	q.Set("state", state)
	q.Set("redirect_uri", s.RedirectURI)

	// Stripe atualmente não obriga PKCE na web, mas enviamos se suportarem no futuro
	authURL.RawQuery = q.Encode()
	return authURL.String()
}

// ExchangeCode troca o código de autorização pelo Access Token da conta conectada (Stripe Account)
func (s *StripeProvider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	tokenURL := "https://connect.stripe.com/oauth/token"

	payload := map[string]string{
		"client_secret": s.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição ao Stripe: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Stripe: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("stripe retornou erro %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		AccessToken      string `json:"access_token"`
		StripeUserID     string `json:"stripe_user_id"` // ID da conta conectada (acct_...)
		StripePublishKey string `json:"stripe_publishable_key"`
		RefreshToken     string `json:"refresh_token"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do Stripe: %w", err)
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    0, // Tokens do Stripe Connect padrão não expiram automaticamente
		TokenType:    "bearer",
	}, nil
}
