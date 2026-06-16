package swarm

import (
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/swarm/agents"
	"gorm.io/gorm"
)

// Supervisor gerencia o ciclo de vida dos Agentes dinâmicos baseados no Banco de Dados.
type Supervisor struct {
	EventBus *evbus.EventBus
	DB       *gorm.DB
}

func NewSupervisor(bus *evbus.EventBus, db *gorm.DB) *Supervisor {
	return &Supervisor{
		EventBus: bus,
		DB:       db,
	}
}

// Start lê todos os agentes do banco e sobe uma goroutine para cada um.
func (s *Supervisor) Start() error {
	log.Println("🌟 [Supervisor] Iniciando Motor Dinâmico de Agentes (The Infinite Swarm)...")

	var activeAgents []models.Agent
	if err := s.DB.Where("type = ?", "llm").Find(&activeAgents).Error; err != nil {
		return err
	}

	log.Printf("🌟 [Supervisor] %d Agentes encontrados no banco de dados. Dando boot...", len(activeAgents))

	for _, dbAgent := range activeAgents {
		s.SpawnAgent(dbAgent)
	}

	// TODO: Ouvir tópico NATS `system.agent.spawn` para subir novos agentes recém-cadastrados sem reboot
	return nil
}

// SpawnAgent instancia e inicializa um GenericAgent
func (s *Supervisor) SpawnAgent(agentModel models.Agent) {
	genericAgent := agents.NewGenericAgent(s.EventBus, s.DB, agentModel)
	
	// Executa em uma goroutine separada para não bloquear o boot
	go func(a *agents.GenericAgent) {
		if err := a.Start(); err != nil {
			log.Printf("❌ [Supervisor] Falha ao iniciar agente '%s': %v", a.Model.Name, err)
		}
	}(genericAgent)
}
