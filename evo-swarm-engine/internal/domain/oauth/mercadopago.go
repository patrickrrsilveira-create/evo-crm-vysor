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

// MercadoPagoProvider implementa a integração OAuth para o Mercado Pago
type MercadoPagoProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func NewMercadoPagoProvider() *MercadoPagoProvider {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}
	redirectURI := fmt.Sprintf("%s/api/v1/oauth/callback/mercadopago", backendURL)

	return &MercadoPagoProvider{
		ClientID:     os.Getenv("MERCADOPAGO_CLIENT_ID"),
		ClientSecret: os.Getenv("MERCADOPAGO_CLIENT_SECRET"),
		RedirectURI:  redirectURI,
	}
}

func (m *MercadoPagoProvider) Name() string {
	return "mercadopago"
}

// BuildAuthURL constrói a URL de login do Mercado Pago
func (m *MercadoPagoProvider) BuildAuthURL(state, codeChallenge string) string {
	authURL, _ := url.Parse("https://auth.mercadopago.com/authorization")
	q := authURL.Query()
	q.Set("client_id", m.ClientID)
	q.Set("response_type", "code")
	q.Set("platform_id", "mp")
	q.Set("state", state)
	q.Set("redirect_uri", m.RedirectURI)

	authURL.RawQuery = q.Encode()
	return authURL.String()
}

// ExchangeCode troca o código de autorização pelo Access Token do vendedor
func (m *MercadoPagoProvider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	tokenURL := "https://api.mercadopago.com/oauth/token"

	data := url.Values{}
	data.Set("client_secret", m.ClientSecret)
	data.Set("client_id", m.ClientID)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", m.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição ao Mercado Pago: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Mercado Pago: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("mercadopago retornou erro %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		PublicKey    string `json:"public_key"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do Mercado Pago: %w", err)
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "bearer",
	}, nil
}
