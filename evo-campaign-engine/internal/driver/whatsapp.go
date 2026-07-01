package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type WhatsAppDriver struct {
	evolutionURL string
	n8nURL       string
	httpClient   *http.Client
}

func NewWhatsApp(evolutionURL, n8nURL string) *WhatsAppDriver {
	return &WhatsAppDriver{
		evolutionURL: evolutionURL,
		n8nURL:       n8nURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (d *WhatsAppDriver) Send(ctx context.Context, instance, recipient string, content Content) SendResult {
	if content.MediaURL != "" {
		return d.sendMedia(ctx, instance, recipient, content)
	}
	return d.sendText(ctx, instance, recipient, content.Text)
}

func (d *WhatsAppDriver) sendText(ctx context.Context, instance, recipient, text string) SendResult {
	if d.evolutionURL == "" {
		return SendResult{OK: false, Error: "evolution_api_url not configured"}
	}

	url := fmt.Sprintf("%s/message/sendText/%s", d.evolutionURL, instance)
	payload := map[string]interface{}{
		"number": recipient,
		"text":   text,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return SendResult{OK: false, Error: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return SendResult{OK: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("whatsapp text sent to %s via %s (status=%d)", recipient, instance, resp.StatusCode)
		return SendResult{OK: true, ProviderStatus: fmt.Sprintf("http_%d", resp.StatusCode)}
	}

	return SendResult{
		OK:             false,
		ProviderStatus: fmt.Sprintf("http_%d", resp.StatusCode),
		Error:          fmt.Sprintf("evolution-api returned %d", resp.StatusCode),
	}
}

func (d *WhatsAppDriver) sendMedia(ctx context.Context, instance, recipient string, content Content) SendResult {
	if d.n8nURL == "" {
		return SendResult{OK: false, Error: "n8n_webhook_url not configured"}
	}

	payload := map[string]interface{}{
		"telefone":  recipient,
		"video_url": content.MediaURL,
		"media":     content.MediaURL,
		"instance":  instance,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.n8nURL, bytes.NewReader(body))
	if err != nil {
		return SendResult{OK: false, Error: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return SendResult{OK: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("whatsapp media sent to %s via n8n (instance=%s status=%d)", recipient, instance, resp.StatusCode)
		return SendResult{OK: true, ProviderStatus: fmt.Sprintf("http_%d", resp.StatusCode)}
	}

	return SendResult{
		OK:             false,
		ProviderStatus: fmt.Sprintf("http_%d", resp.StatusCode),
		Error:          fmt.Sprintf("n8n webhook returned %d", resp.StatusCode),
	}
}
