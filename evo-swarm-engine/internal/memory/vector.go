package memory

import (
	"context"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database"
	"github.com/pgvector/pgvector-go"
)

// DocumentMemory representa a tabela de memórias vetoriais do Agente
type DocumentMemory struct {
	ID        int64           `gorm:"primaryKey;autoIncrement"`
	AgentID   string          `gorm:"index"` // ID do agente ou workflow dono da memória
	Content   string          `gorm:"type:text"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)"` // Padrão OpenAI (text-embedding-ada-002 ou text-embedding-3-small)
	Metadata  string          `gorm:"type:jsonb"`        // Metadados adicionais em JSON (URL, título, etc)
}

func (DocumentMemory) TableName() string {
	return "agent_memories"
}

// MemoryEngine gerencia o armazenamento e recuperação de contexto por semântica
type MemoryEngine struct{}

func NewMemoryEngine() *MemoryEngine {
	// Auto-Migra a tabela de memórias (Garante que a extensão vetor esteja ativa no Postgres)
	database.DB.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err := database.DB.AutoMigrate(&DocumentMemory{}); err != nil {
		log.Printf("Aviso: Falha ao migrar tabela de memórias: %v", err)
	}

	return &MemoryEngine{}
}

// StoreMemory salva uma nova memória (texto + vetor) para um agente
func (m *MemoryEngine) StoreMemory(ctx context.Context, agentID, content string, embedding []float32, metadata string) error {
	doc := DocumentMemory{
		AgentID:   agentID,
		Content:   content,
		Embedding: pgvector.NewVector(embedding),
		Metadata:  metadata,
	}
	return database.DB.Create(&doc).Error
}

// SearchSimilar busca as N memórias mais próximas (Busca Semântica usando Cosine Similarity)
func (m *MemoryEngine) SearchSimilar(ctx context.Context, agentID string, queryEmbedding []float32, limit int) ([]DocumentMemory, error) {
	var results []DocumentMemory
	vec := pgvector.NewVector(queryEmbedding)

	// L2 distance (<->) ou Cosine Similarity (<=>). Vamos usar Cosine Similarity:
	err := database.DB.
		Where("agent_id = ?", agentID).
		Order("embedding <=> ?").
		Limit(limit).
		Find(&results, vec).Error

	return results, err
}
