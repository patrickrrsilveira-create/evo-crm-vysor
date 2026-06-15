package main

import (
	"context"
	"log"
	"os"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/adapters"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/coordinator"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/mcp"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/memory"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/middleware"
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
	database.Connect()

	// Conecta ao NATS (Event Bus)
	events.ConnectNATS()
	defer events.GlobalEventBus.Close()

	// Conecta ao Redis (Memory Engine)
	database.ConnectRedis()

	// Inicializa o Swarm Coordinator
	coord := coordinator.NewCoordinator(events.GlobalEventBus)
	if err := coord.Start(); err != nil {
		log.Fatalf("Falha ao iniciar o Coordinator: %v", err)
	}

	// Inicializa os Workers do Swarm
	agentWorker := workers.NewAgentWorker(events.GlobalEventBus)
	if err := agentWorker.Start(); err != nil {
		log.Fatalf("Falha ao iniciar o Agent Worker: %v", err)
	}

	proactiveEngine := workers.NewProactiveEngine(events.GlobalEventBus)
	proactiveEngine.Start()

	// Inicializa o Client MCP (Conexão com as Ferramentas)
	mcpClient := mcp.NewClient("http://localhost:3000/sse") // Exemplo de servidor MCP externo
	if err := mcpClient.Connect(context.Background()); err != nil {
		log.Printf("Aviso: Falha ao conectar ao servidor MCP: %v", err)
	} else {
		mcpClient.DiscoverTools(context.Background())
	}

	// Inicializa Adaptadores Ativos (que não dependem de rota HTTP)

	adapters.NewTeamsAdapter(events.GlobalEventBus, "appid", "pass").Start(context.Background())
	adapters.NewEmailAdapter("smtp", "imap", "user", "pass").Start(context.Background())
	adapters.NewCalendarAdapter("creds.json").Start(context.Background())

	// Nota: Google Drive não implementa Start() de escuta passiva, apenas métodos ativos
	adapters.NewGoogleDriveAdapter("creds.json")

	// Inicializa a Memory Engine (Vector Database via PostgreSQL/PGVector)
	memory.NewMemoryEngine()

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
	app.Use("/api", middleware.EvoAuthMiddleware())

	// Adapters e Webhooks Omni-Channel
	evolutionAdapter := adapters.NewEvolutionAdapter(events.GlobalEventBus, "http://evo:8080", "123")
	chatwootAdapter := adapters.NewChatwootMirrorAdapter(events.GlobalEventBus)
	a2aAdapter := adapters.NewA2AAdapter(events.GlobalEventBus)

	evolutionAdapter.Start(context.Background())
	chatwootAdapter.Start(context.Background())

	// Registra rotas após middleware
	evolutionAdapter.RegisterWebhookRoute(app)
	chatwootAdapter.RegisterWebhookRoute(app)
	a2aAdapter.RegisterRoutes(app)

	// Rotas do Kubernetes (Probes)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		// Checar Banco e Redis
		sqlDB, err := database.DB.DB()
		dbStatus := "ok"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "error"
		}

		redisStatus := "ok"
		if database.RedisClient.Ping(context.Background()).Err() != nil {
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

	// Pega a porta do .env ou usa 8001 (para rodar ao lado do Python na 8000)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("🚀 Swarm Engine rodando na porta %s...", port)
	log.Fatal(app.Listen(":" + port))
}
