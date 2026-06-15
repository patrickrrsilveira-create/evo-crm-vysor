package main

import (
	"log"
	"os"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Conecta ao PostgreSQL
	database.Connect()

	// Inicializa o Fiber (Framework Web Ultra-rápido)
	app := fiber.New(fiber.Config{
		AppName: "Evo Swarm Engine",
	})

	// Rota de Healthcheck básica
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "online",
			"engine": "Go Swarm",
			"version": "1.0.0",
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
