package crm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client é o conector nativo em Go para a API Core do Evo CRM (Ruby on Rails).
// Ele espelha o comportamento do `EvoCrmClient` em Python (src/services/adk/tools/evo_crm/base.py)
type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient inicializa o conector lendo as variáveis de ambiente necessárias.
func NewClient() *Client {
	baseURL := os.Getenv("EVO_AI_CRM_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	apiToken := os.Getenv("EVOAI_CRM_API_TOKEN")

	return &Client{
		BaseURL:  baseURL,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest encapsula a construção e disparo da requisição HTTP com os cabeçalhos de segurança.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) ([]byte, error) {
	if c.APIToken == "" {
		return nil, fmt.Errorf("EVOAI_CRM_API_TOKEN não está configurado no ambiente")
	}

	// Tratamento do endpoint para sempre iniciar com /api/v1/
	cleanEndpoint := endpoint
	if !strings.HasPrefix(cleanEndpoint, "/api/v1/") {
		if strings.HasPrefix(cleanEndpoint, "/") {
			cleanEndpoint = "/api/v1" + cleanEndpoint
		} else {
			cleanEndpoint = "/api/v1/" + cleanEndpoint
		}
	}

	fullURL := c.BaseURL + cleanEndpoint

	var bodyReader io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("falha ao encodar payload JSON: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição HTTP: %w", err)
	}

	// Cabeçalhos de Segurança da Arquitetura Core
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Service-Token", c.APIToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro de rede na chamada ao CRM: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta da API do CRM: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API do CRM retornou erro HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}

// Post executa uma requisição HTTP POST genérica no CRM
func (c *Client) Post(ctx context.Context, endpoint string, data interface{}) ([]byte, error) {
	return c.makeRequest(ctx, http.MethodPost, endpoint, data)
}

// Get executa uma requisição HTTP GET genérica no CRM
func (c *Client) Get(ctx context.Context, endpoint string) ([]byte, error) {
	return c.makeRequest(ctx, http.MethodGet, endpoint, nil)
}
