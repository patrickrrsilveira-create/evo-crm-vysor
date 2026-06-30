package oauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// NotionProvider implementa a interface Provider para a API V2 do Notion.
type NotionProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// NewNotionProvider constrói a estratégia OAuth para o Notion
func NewNotionProvider() *NotionProvider {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}
	redirectURI := fmt.Sprintf("%s/api/v1/oauth/callback/notion", backendURL)

	return &NotionProvider{
		ClientID:     os.Getenv("NOTION_CLIENT_ID"),
		ClientSecret: os.Getenv("NOTION_CLIENT_SECRET"),
		RedirectURI:  redirectURI,
	}
}

// Name retorna o nome da variante "notion"
func (n *NotionProvider) Name() string {
	return "notion"
}

// BuildAuthURL constrói a URL de autorização pública do Notion
func (n *NotionProvider) BuildAuthURL(state, codeChallenge string) string {
	authURL, _ := url.Parse("https://api.notion.com/v1/oauth/authorize")
	q := authURL.Query()
	q.Set("client_id", n.ClientID)
	q.Set("response_type", "code")
	q.Set("owner", "user")
	q.Set("redirect_uri", n.RedirectURI)
	q.Set("state", state)

	// O Notion na V2 não documenta PKCE explicitamente, mas enviamos para consistência de segurança.
	// Se o Notion rejeitar, removeremos o PKCE.
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")

	authURL.RawQuery = q.Encode()
	return authURL.String()
}

// ExchangeCode troca o código pelo Access Token
// O Notion exige HTTP Basic Auth no Header usando Base64(client_id:client_secret)
func (n *NotionProvider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	tokenURL := "https://api.notion.com/v1/oauth/token"

	payload := map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": n.RedirectURI,
		// "code_verifier": codeVerifier, // Notion doesn't officially support PKCE validation on token exchange
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// Notion exige Basic Auth
	authString := fmt.Sprintf("%s:%s", n.ClientID, n.ClientSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição de token ao Notion: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Notion: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("notion retornou erro %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		AccessToken   string `json:"access_token"`
		WorkspaceName string `json:"workspace_name"`
		WorkspaceIcon string `json:"workspace_icon"`
		WorkspaceID   string `json:"workspace_id"`
		BotID         string `json:"bot_id"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON de tokens do notion: %w", err)
	}

	// Notion Access Tokens não expiram por padrão (ExpiresIn = 0), eles duram até o usuário revogar.
	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: "", // Notion não usa Refresh Token clássico na V2 pública
		ExpiresIn:    0,
		TokenType:    "bearer",
	}, nil
}
