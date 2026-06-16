package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/nats-io/nats.go"
)

// Capability define as habilidades de um Agente
type Capability struct {
	AgentID     string   `json:"agent_id"`
	Subject     string   `json:"subject"` // Ex: agent.crm.task
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Registry gerencia a listagem dinâmica de Agentes no Swarm
type Registry struct {
	KV nats.KeyValue
}

func NewRegistry(bus *evbus.EventBus) (*Registry, error) {
	js, err := bus.Conn.JetStream()
	if err != nil {
		return nil, err
	}

	kv, err := js.KeyValue("CAPABILITIES")
	if err != nil {
		kv, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      "CAPABILITIES",
			Description: "Registro Dinâmico de Agentes do Swarm",
			TTL:         10 * time.Minute, // Agentes devem fazer heartbeat/re-register a cada 5 min
		})
		if err != nil {
			return nil, fmt.Errorf("falha ao criar KV CAPABILITIES: %v", err)
		}
	}

	return &Registry{KV: kv}, nil
}

// Register publica as capacidades de um Agente no Bucket
func (r *Registry) Register(cap Capability) error {
	cap.UpdatedAt = time.Now()
	data, err := json.Marshal(cap)
	if err != nil {
		return fmt.Errorf("falha ao serializar capacidade: %v", err)
	}

	_, err = r.KV.Put(cap.AgentID, data)
	if err == nil {
		log.Printf("🛂 [Registry] Agente Registrado: %s -> %s", cap.AgentID, cap.Subject)
	}
	return err
}

// GetAllCapabilities retorna todos os agentes atualmente ativos
func (r *Registry) GetAllCapabilities() ([]Capability, error) {
	keys, err := r.KV.Keys()
	if err != nil {
		if err == nats.ErrKeyNotFound {
			return []Capability{}, nil
		}
		return nil, err
	}

	var caps []Capability
	for _, k := range keys {
		entry, err := r.KV.Get(k)
		if err == nil {
			var c Capability
			if err := json.Unmarshal(entry.Value(), &c); err != nil {
				log.Printf("⚠️ [Registry] Erro ao decodificar capacidade para chave %s: %v", k, err)
				continue
			}
			caps = append(caps, c)
		}
	}
	return caps, nil
}
