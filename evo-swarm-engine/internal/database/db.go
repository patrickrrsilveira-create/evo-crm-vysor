package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB é a instância global de conexão com o banco
var DB *gorm.DB

// Connect inicializa a conexão com o PostgreSQL
func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Falha: DATABASE_URL não definida no arquivo .env")
	}

	// Conecta usando GORM com log no modo Silent para não sujar o terminal em dev
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("Falha crítica: Não foi possível conectar ao banco de dados: %v", err)
	}

	log.Println("✅ Conectado ao PostgreSQL com sucesso!")

	DB = db
}
