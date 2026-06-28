package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
)

// RegisterFinanceRoutes registra os webhooks para Stripe e Mercado Pago
func RegisterFinanceRoutes(app *fiber.App, eventBus *evbus.EventBus) {
	webhooks := app.Group("/api/v1/webhooks/finance")

	webhooks.Post("/stripe", func(c *fiber.Ctx) error {
		// Stripe exige validação da assinatura baseada no Raw Body
		signature := c.Get("Stripe-Signature")
		if signature == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Stripe-Signature"})
		}

		webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		if webhookSecret != "" {
			rawBody := c.Body()
			if !verifyStripeSignature(rawBody, signature, webhookSecret) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Signature"})
			}
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(c.Body(), &payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		// Despacha evento de pagamento no NATS para a IA ser notificada
		eventType, _ := payload["type"].(string)
		if eventType == "checkout.session.completed" || eventType == "payment_intent.succeeded" {
			dataBytes, _ := json.Marshal(payload)
			eventBus.Publish("finance.payment.approved", dataBytes)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	webhooks.Post("/mercadopago", func(c *fiber.Ctx) error {
		// MP validation requires reading the x-signature and x-request-id headers
		signature := c.Get("x-signature")
		requestID := c.Get("x-request-id")

		webhookSecret := os.Getenv("MERCADOPAGO_WEBHOOK_SECRET")
		if webhookSecret != "" && signature != "" && requestID != "" {
			action := c.Query("data.id")
			if !verifyMercadoPagoSignature(signature, requestID, action, webhookSecret) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Signature"})
			}
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(c.Body(), &payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		actionType, _ := payload["action"].(string)
		if actionType == "payment.created" || actionType == "payment.updated" {
			// Na vida real, precisaríamos buscar o status do pagamento na API usando o ID recebido
			// Aqui apenas despachamos o aviso
			dataBytes, _ := json.Marshal(payload)
			eventBus.Publish("finance.payment.mercadopago", dataBytes)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

// verifyStripeSignature valida a assinatura do webhook do Stripe
func verifyStripeSignature(payload []byte, header string, secret string) bool {
	var timestamp string
	var sig string

	parts := strings.Split(header, ",")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			if kv[0] == "t" {
				timestamp = kv[1]
			} else if kv[0] == "v1" {
				sig = kv[1]
			}
		}
	}

	if timestamp == "" || sig == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%s.%s", timestamp, string(payload))))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expectedSig))
}

// verifyMercadoPagoSignature valida a assinatura do webhook do Mercado Pago (v2)
func verifyMercadoPagoSignature(xSignature, xRequestID, dataID, secret string) bool {
	// A assinatura do MP no header x-signature vem no formato: ts=...,v1=...
	var timestamp string
	var sig string

	parts := strings.Split(xSignature, ",")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			if kv[0] == "ts" {
				timestamp = kv[1]
			} else if kv[0] == "v1" {
				sig = kv[1]
			}
		}
	}

	if timestamp == "" || sig == "" {
		return false
	}

	// O payload esperado pelo MP para HMAC é: id,request-id,ts
	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, xRequestID, timestamp)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(manifest))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expectedSig))
}
