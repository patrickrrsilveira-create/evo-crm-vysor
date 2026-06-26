package workers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ProactiveCampaign struct {
	ID              int64  `gorm:"primaryKey"`
	AccountID       int64  `gorm:"column:account_id"`
	TriggerTarget   string `gorm:"column:trigger_target"`
	DelayHours      int    `gorm:"column:delay_hours"`
	MessageTemplate string `gorm:"column:message_template"`
	AgentID         string `gorm:"column:agent_id"`
	Status          string `gorm:"column:status"`
}

type ProactiveEngine struct {
	Cron     *cron.Cron
	EventBus *evbus.EventBus
	db       *gorm.DB
}

func NewProactiveEngine(bus *evbus.EventBus, db *gorm.DB) *ProactiveEngine {
	// Cria um Cron job logger default
	c := cron.New()
	return &ProactiveEngine{
		Cron:     c,
		EventBus: bus,
		db:       db,
	}
}

func (p *ProactiveEngine) Start() {
	log.Println("⏰ [ProactiveEngine] Inicializando motor de Campanhas Proativas (Cron)...")

	// Roda a cada 60 minutos: "0 * * * *"
	_, err := p.Cron.AddFunc("@hourly", p.runCampaignsLoop)
	if err != nil {
		log.Fatalf("Erro ao iniciar Cron job: %v", err)
	}

	p.Cron.Start()
	log.Println("⏰ [ProactiveEngine] Cron rodando em background")
}

func (p *ProactiveEngine) runCampaignsLoop() {
	log.Println("🔍 [ProactiveEngine] Varrendo banco de dados por Campanhas Ativas...")

	var campaigns []ProactiveCampaign
	// raw sql query similar to Python: SELECT id, account_id, trigger_target, delay_hours, message_template, agent_id FROM proactive_campaigns WHERE status = 'ACTIVE'
	// Assume that the table exists or will exist in core database
	result := p.db.Table("proactive_campaigns").Where("status = ?", "ACTIVE").Find(&campaigns)

	if result.Error != nil {
		log.Printf("⚠️ [ProactiveEngine] Tabela proactive_campaigns não encontrada ou erro: %v", result.Error)
		return
	}

	for _, camp := range campaigns {
		log.Printf("🚀 [ProactiveEngine] Processando Campanha %d (Target: %s)", camp.ID, camp.TriggerTarget)

		// Simula disparo escalonado (no Python era await asyncio.sleep(30))
		// Aqui lançamos uma goroutine que fará o delay se necessário
		go p.queueDelivery(camp)
	}
}

func (p *ProactiveEngine) queueDelivery(camp ProactiveCampaign) {
	// Lógica simplificada de fila de disparo anti-ban:
	// A engine apenas injeta EventMessageReceived no NATS com as tags ocultas
	// indicando que foi um trigger interno (proativo).

	agentUUID, _ := uuid.Parse(camp.AgentID)

	startEvent := events.AgentStartedEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New(),
			EventType: events.EventAgentStarted,
			Timestamp: time.Now(),
		},
		AgentID: agentUUID,
		TaskID:  uuid.New(),
		Payload: "TRIGGER_PROACTIVE_CAMPAIGN: " + camp.MessageTemplate,
	}

	eventData, _ := json.Marshal(startEvent)
	p.EventBus.Publish(string(events.EventAgentStarted), eventData)

	log.Printf("✉️ [ProactiveEngine] Campanha %d injetada no barramento para o Agente %s", camp.ID, camp.AgentID)
}
