package knowledge

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"gorm.io/gorm"
)

// LegacyPGVectorStore implementa vectorstores.VectorStore do Langchaingo
// Mas adaptado para usar o nosso schema legado (knowledge_document_chunks)
type LegacyPGVectorStore struct {
	db       *gorm.DB
	embedder embeddings.Embedder
}

var _ vectorstores.VectorStore = &LegacyPGVectorStore{}

// NewLegacyPGVectorStore cria uma nova instância
func NewLegacyPGVectorStore(db *gorm.DB, embedder embeddings.Embedder) *LegacyPGVectorStore {
	return &LegacyPGVectorStore{
		db:       db,
		embedder: embedder,
	}
}

func (s *LegacyPGVectorStore) AddDocuments(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error) {
	// O Upsert de documentos não é o foco do Runtime agora (é feito pelo Frontend/Python),
	// mas deixamos um stub para cumprir a interface.
	return nil, fmt.Errorf("AddDocuments not implemented in legacy store yet")
}

func (s *LegacyPGVectorStore) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
	opts := vectorstores.Options{}
	for _, opt := range options {
		opt(&opts)
	}

	// 1. Gera embedding usando o LangChain Embedder
	queryVectors, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar embedding via langchaingo: %w", err)
	}
	
	if len(queryVectors) == 0 {
		return nil, nil
	}

	queryEmbedding := pgvector.NewVector(queryVectors)

	// 2. Extrai filtros do metadata
	var kbIDs []string
	if opts.Filters != nil {
		if filters, ok := opts.Filters.(map[string]any); ok {
			if kbIDList, ok := filters["knowledge_base_ids"].([]string); ok {
				kbIDs = kbIDList
			}
		}
	}

	if len(kbIDs) == 0 {
		return nil, nil // Sem permissões
	}

	// 3. Consulta vetorial com o GORM
	type Result struct {
		Content  string
		Distance float32
	}
	var results []Result

	sqlQuery := `
		SELECT c.content, c.embedding <=> ? AS distance
		FROM knowledge_document_chunks c
		JOIN knowledge_documents d ON c.knowledge_document_id = d.id
		WHERE d.knowledge_base_id IN ?
		ORDER BY c.embedding <=> ?
		LIMIT ?
	`

	err = s.db.Raw(sqlQuery, queryEmbedding, kbIDs, queryEmbedding, numDocuments).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar contexto vetorial no banco: %w", err)
	}

	// 4. Mapeia para schema.Document do Langchaingo
	docs := make([]schema.Document, 0, len(results))
	for _, res := range results {
		docs = append(docs, schema.Document{
			PageContent: res.Content,
			Score:       res.Distance,
			Metadata:    map[string]any{"source": "legacy_db"},
		})
	}

	return docs, nil
}
