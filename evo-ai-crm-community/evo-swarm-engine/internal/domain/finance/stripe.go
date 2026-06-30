package finance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// StripeProvider implementa as ações financeiras para a Stripe
type StripeProvider struct {
	SecretKey string
}

func NewStripeProvider() *StripeProvider {
	return &StripeProvider{
		SecretKey: os.Getenv("STRIPE_SECRET_KEY"), // Na arquitetura final, isso virá do banco de dados (OAuth do usuário)
	}
}

// CreatePaymentLink cria um link de pagamento na Stripe via API REST
func (s *StripeProvider) CreatePaymentLink(amount float64, currency string, description string) (string, error) {
	if s.SecretKey == "" {
		return "", fmt.Errorf("STRIPE_SECRET_KEY não configurada")
	}

	// 1. Criar um Produto/Preço ou usar um Preço ad-hoc na Stripe (Price API)
	// Para links de pagamento ad-hoc na v1/payment_links precisamos de um PriceID,
	// ou podemos usar a Checkout Session API para cobranças dinâmicas.
	// Vamos usar Checkout Session para permitir valor flexível.

	url := "https://api.stripe.com/v1/checkout/sessions"

	// Stripe espera valores em centavos
	amountCents := int(amount * 100)

	payload := fmt.Sprintf(
		"payment_method_types[]=card&line_items[0][price_data][currency]=%s&line_items[0][price_data][product_data][name]=%s&line_items[0][price_data][unit_amount]=%d&line_items[0][quantity]=1&mode=payment&success_url=https://evo.com/success",
		currency, description, amountCents,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+s.SecretKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("stripe api error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.URL, nil
}
