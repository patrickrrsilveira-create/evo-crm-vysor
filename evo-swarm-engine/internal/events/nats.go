package events

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// EventBus encapsula a conexão do NATS para publish/subscribe.
type EventBus struct {
	Conn *nats.Conn
}

// GlobalEventBus é a instância global de mensageria.
var GlobalEventBus *EventBus

// ConnectNATS inicializa a conexão com o servidor NATS.
func ConnectNATS() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	// Tenta conectar ao NATS com retentativas automáticas
	nc, err := nats.Connect(natsURL, nats.RetryOnFailedConnect(true), nats.MaxReconnects(5), nats.ReconnectWait(2*time.Second))
	if err != nil {
		log.Fatalf("Falha crítica: Não foi possível conectar ao NATS (%s): %v", natsURL, err)
	}

	log.Printf("✅ Conectado ao NATS com sucesso! (URL: %s)", natsURL)

	GlobalEventBus = &EventBus{
		Conn: nc,
	}
}

// Publish publica um evento no NATS de forma assíncrona.
func (eb *EventBus) Publish(subject string, data []byte) error {
	return eb.Conn.Publish(subject, data)
}

// Subscribe assina um tópico (subject) no NATS.
func (eb *EventBus) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return eb.Conn.Subscribe(subject, handler)
}

// Close encerra a conexão de forma limpa.
func (eb *EventBus) Close() {
	if eb.Conn != nil {
		eb.Conn.Close()
		log.Println("Conexão com NATS encerrada.")
	}
}
