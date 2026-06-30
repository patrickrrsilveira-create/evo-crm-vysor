package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/finance"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/nats-io/nats.go"
)

// FinanceBridgeAdapter escuta comandos NATS do FinanceAgent e chama a Stripe/MercadoPago
type FinanceBridgeAdapter struct {
	EventBus *evbus.EventBus
}

func NewFinanceBridgeAdapter(bus *evbus.EventBus) *FinanceBridgeAdapter {
	return &FinanceBridgeAdapter{
		EventBus: bus,
	}
}

func (a *FinanceBridgeAdapter) Start(ctx context.Context) error {
	log.Println("💰 [FinanceBridge] Iniciado. Traduzindo comandos NATS -> Gateways de Pagamento")

	_, err := a.EventBus.Conn.QueueSubscribe("finance.command.charge", "finance_bridge_group", a.handleChargeCommand)
	if err != nil {
		return err
	}

	return nil
}

func (a *FinanceBridgeAdapter) handleChargeCommand(msg *nats.Msg) {
	var args struct {
		Provider    string  `json:"provider"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		Description string  `json:"description"`
	}

	if err := json.Unmarshal(msg.Data, &args); err != nil {
		log.Printf("❌ [FinanceBridge] Erro ao decodificar comando: %v", err)
		msg.Respond([]byte(`{"error": "invalid payload"}`))
		return
	}

	log.Printf("💰 [FinanceBridge] Solicitando link via %s (%.2f %s)", args.Provider, args.Amount, args.Currency)

	// Aqui invocaríamos o Provider real da Stripe que está em domain/finance ou domain/oauth
	// Exemplo simulado da integração real
	
	// Como o adapter real depende do Token da Empresa logada, vamos abstrair a geração
	stripeProvider := finance.NewStripeProvider()
	
	// Chamada real mockada no provider interno
	url, err := stripeProvider.CreatePaymentLink(args.Amount, args.Currency, args.Description)

	if err != nil {
		log.Printf("❌ [FinanceBridge] Erro na Stripe: %v", err)
		msg.Respond([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
		return
	}

	// Responde com Sucesso para o FinanceAgent
	type FinanceResponse struct {
		Status string `json:"status"`
		URL    string `json:"url"`
	}

	responsePayload, err := json.Marshal(FinanceResponse{
		Status: "success",
		URL:    url,
	})
	if err != nil {
		log.Printf("❌ [FinanceBridge] Erro ao serializar FinanceResponse: %v", err)
		return
	}
	
	if msg.Reply != "" {
		if err := msg.Respond(responsePayload); err != nil {
			log.Printf("❌ [FinanceBridge] Erro ao responder via NATS: %v", err)
		}
	}
}
