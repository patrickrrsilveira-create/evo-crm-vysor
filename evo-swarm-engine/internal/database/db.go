package database

import (
	"fmt"
	"log"
	"os"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect inicializa a conexão com o PostgreSQL
func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL não definida no arquivo .env")
	}

	// Conecta usando GORM com log no modo Silent para não sujar o terminal em dev
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados: %w", err)
	}

	log.Println("✅ Conectado ao PostgreSQL com sucesso!")

	return db, nil
}

// AutoMigrate roda a migração dos modelos
func AutoMigrate(db *gorm.DB) error {
	log.Println("Rodando AutoMigrate...")
	err := db.AutoMigrate(
		&models.OAuthSession{},
		&models.ConversationMessage{},
	)
	if err != nil {
		return fmt.Errorf("falha no AutoMigrate: %w", err)
	}
	return nil
}
