package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/adapters"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/api/routes"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/coordinator"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/oauth"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/memory"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/middleware"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/services"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/swarm"
	swarmservices "github.com/PatrickRSilveira/evo-swarm-engine/internal/swarm/services"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/workers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"time"
)

func main() {
	// Carrega as variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Conecta ao PostgreSQL
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Falha crítica ao conectar ao PostgreSQL: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Printf("Aviso: Falha no AutoMigrate: %v", err)
	}

	// Conecta ao NATS (Event Bus)
	evbus.ConnectNATS()
	defer evbus.GlobalEventBus.Close()

	// Conecta ao Redis (Memory Engine)
	redisClient, err := database.ConnectRedis()
	if err != nil {
		log.Fatalf("Falha crítica ao conectar ao Redis: %v", err)
	}

	// Inicializa o Event Store (Event Sourcing Audit)
	eventStore, err := evbus.NewEventStore(evbus.GlobalEventBus)
	if err != nil {
		log.Printf("⚠️ Aviso: Falha ao inicializar EventStore: %v", err)
	} else {
		// Log de teste para auditoria
		eventStore.LogEvent("system.events.started", "system", "engine_boot", map[string]string{"version": "1.0"})
	}

	// Inicializa o Swarm Coordinator
	coord := coordinator.NewCoordinator(evbus.GlobalEventBus, db)
	if err := coord.Start(); err != nil {
		log.Fatalf("Falha ao iniciar o Coordinator: %v", err)
	}

	// Inicializa Microsserviços NATS (Desacoplados)
	memService := swarmservices.NewMemoryService(evbus.GlobalEventBus, db, redisClient)
	if err := memService.Start(); err != nil {
		log.Fatal(err)
	}

	ragService := swarmservices.NewRAGService(evbus.GlobalEventBus, db)
	if err := ragService.Start(); err != nil {
		log.Fatal(err)
	}

	// Inicializa o The Infinite Swarm (Supervisor Dinâmico)
	swarmSupervisor := swarm.NewSupervisor(evbus.GlobalEventBus, db)
	if err := swarmSupervisor.Start(); err != nil {
		log.Fatalf("Falha ao iniciar Swarm Supervisor: %v", err)
	}

	proactiveEngine := workers.NewProactiveEngine(evbus.GlobalEventBus, db)
	proactiveEngine.Start()

	// (O antigo setup de MCP Client hardcoded foi removido a favor do Roteamento Dinâmico em mcp_routes.go)

	// Inicializa Adaptadores Ativos (que não dependem de rota HTTP)

	adapters.NewTeamsAdapter(evbus.GlobalEventBus, "appid", "pass").Start(context.Background())
	adapters.NewEmailAdapter("smtp", "imap", "user", "pass").Start(context.Background())
	adapters.NewCalendarAdapter("creds.json").Start(context.Background())

	// Nota: Google Drive não implementa Start() de escuta passiva, apenas métodos ativos
	adapters.NewGoogleDriveAdapter("creds.json")

	// Inicializa a Memory Engine (Vector Database via PostgreSQL/PGVector)
	memory.NewMemoryEngine(db)

	// Inicializa o Fiber (Framework Web Ultra-rápido)
	app := fiber.New(fiber.Config{
		AppName: "Evo Swarm Engine",
	})

	// Setup Middlewares de Segurança e Tráfego
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,             // Máximo de 100 requests por IP
		Expiration: 1 * time.Minute, // Por minuto
	}))

	// EvoAuth Middleware
	app.Use("/api", middleware.EvoAuthMiddleware(db))

	// Adapters e Webhooks Omni-Channel
	evolutionAdapter := adapters.NewEvolutionAdapter(evbus.GlobalEventBus, "http://evo:8080", os.Getenv("EVOLUTION_WEBHOOK_API_KEY"))
	chatwootAdapter := adapters.NewChatwootMirrorAdapter(evbus.GlobalEventBus, db)
	a2aAdapter := adapters.NewA2AAdapter(evbus.GlobalEventBus)
	msTeamsAdapter := adapters.NewMSTeamsAdapter(evbus.GlobalEventBus)
	
	crmBridge := adapters.NewCRMBridgeAdapter(evbus.GlobalEventBus)
	crmBridge.Start(context.Background())

	financeBridge := adapters.NewFinanceBridgeAdapter(evbus.GlobalEventBus)
	financeBridge.Start(context.Background())

	webhookNotifier := adapters.NewWebhookNotifierAdapter(evbus.GlobalEventBus)
	webhookNotifier.Start(context.Background())

	typebotAdapter := adapters.NewTypebotAdapter(evbus.GlobalEventBus)
	typebotAdapter.Start(context.Background())

	n8nAdapter := adapters.NewN8nAdapter(evbus.GlobalEventBus)
	n8nAdapter.Start(context.Background())

	// Rotas OAuth / PKCE
	oauthService := services.NewOAuthService(db)
	oauthService.RegisterProvider(oauth.NewGoogleProvider("google_calendar"))
	oauthService.RegisterProvider(oauth.NewGoogleProvider("google_sheets"))
	oauthService.RegisterProvider(oauth.NewGoogleProvider("google_drive"))
	oauthService.RegisterProvider(oauth.NewNotionProvider())
	oauthService.RegisterProvider(oauth.NewStripeProvider())
	oauthService.RegisterProvider(oauth.NewMercadoPagoProvider())

	routes.RegisterOAuthRoutes(app, oauthService, db)
	routes.RegisterIntegrationRoutes(app)
	routes.RegisterFinanceRoutes(app, evbus.GlobalEventBus)

	// Rotas Dinâmicas de Integração MCP
	routes.RegisterMCPRoutes(app, db)

	// Rotas CRUD da API REST (Fase 1 - Transição do Python)
	routes.RegisterAgentRoutes(app, db)
	routes.RegisterToolRoutes(app, db)

	evolutionAdapter.Start(context.Background())
	chatwootAdapter.Start(context.Background())
	msTeamsAdapter.Start(context.Background())

	// Registra rotas após middleware
	evolutionAdapter.RegisterWebhookRoute(app)
	chatwootAdapter.RegisterWebhookRoute(app)
	a2aAdapter.RegisterRoutes(app, db)
	msTeamsAdapter.RegisterWebhookRoute(app)
	typebotAdapter.RegisterWebhookRoute(app)
	n8nAdapter.RegisterWebhookRoute(app)

	// Rotas do Kubernetes (Probes)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		// Checar Banco e Redis
		sqlDB, err := db.DB()
		dbStatus := "ok"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "error"
		}

		redisStatus := "ok"
		if redisClient.Ping(context.Background()).Err() != nil {
			redisStatus = "error"
		}

		if dbStatus != "ok" || redisStatus != "ok" {
			return c.Status(503).JSON(fiber.Map{
				"status":   "not_ready",
				"database": dbStatus,
				"redis":    redisStatus,
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ready",
			"database": "ok",
			"redis":    "ok",
		})
	})

	// Pega a porta do .env ou usa 8001
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	// Graceful Shutdown Setup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Swarm Engine rodando na porta %s...", port)
		if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor Fiber: %v", err)
		}
	}()

	<-c
	log.Println("Desligando graciosamente...")
	
	// Dá 5 segundos para requisições em andamento terminarem
	_ = app.ShutdownWithTimeout(5 * time.Second)
	
	log.Println("Servidor desligado com sucesso.")
}
