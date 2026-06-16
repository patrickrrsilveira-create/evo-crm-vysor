package knowledge

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"gorm.io/gorm"
)

// RAGService orquestra a Busca Vetorial (Knowledge Nexus).
type RAGService struct {
	db *gorm.DB
}

func NewRAGService(db *gorm.DB) *RAGService {
	return &RAGService{
		db: db,
	}
}

// RetrieveContext busca nas bases de conhecimento os textos mais semelhantes semanticamente usando LangChain Go.
func (s *RAGService) RetrieveContext(ctx context.Context, agentBotID uuid.UUID, apiKey models.APIKey, query string) (string, error) {
	if apiKey.Key == "" {
		return "", fmt.Errorf("openai_api_key ausente")
	}

	// 1. Descobrir quais Bases de Conhecimento estão amarradas a este agente
	var kbLinks []models.KnowledgeBaseAgentBot
	if err := s.db.Where("agent_bot_id = ?", agentBotID).Find(&kbLinks).Error; err != nil {
		return "", fmt.Errorf("falha ao buscar bases atreladas ao agente: %w", err)
	}

	if len(kbLinks) == 0 {
		return "", nil // O agente não tem KBs, retorna vazio silenciosamente.
	}

	kbIDs := make([]string, len(kbLinks))
	for i, link := range kbLinks {
		kbIDs[i] = link.KnowledgeBaseID
	}

	// 2. Inicializa o Client Oficial do LangChain para a OpenAI
	llm, err := openai.New(
		openai.WithToken(apiKey.Key),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	)
	if err != nil {
		return "", fmt.Errorf("falha ao instanciar langchain openai client: %w", err)
	}
	
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return "", fmt.Errorf("falha ao instanciar langchain embedder: %w", err)
	}

	// 3. Instancia o VectorStore Adapter do LangChain
	// (Nosso adapter faz ponte com as tabelas legadas do PostgreSQL)
	store := NewLegacyPGVectorStore(s.db, embedder)

	// 4. Cria o Retriever oficial do LangChain, injetando os filtros por Metadata (Isolamento do Agente)
	retriever := vectorstores.ToRetriever(store, 5, vectorstores.WithFilters(map[string]any{
		"knowledge_base_ids": kbIDs,
	}))

	// 5. Executa a Busca Relevante via Langchain
	docs, err := retriever.GetRelevantDocuments(ctx, query)
	if err != nil {
		return "", fmt.Errorf("falha ao executar retriever do langchain: %w", err)
	}

	if len(docs) == 0 {
		return "", nil
	}

	log.Printf("🧠 [LangChain RAG] Encontrados %d chunks para a query '%s'", len(docs), query)

	// 6. Concatena o texto para ser injetado no System Prompt
	return FormatDocuments(docs), nil
}

// FormatDocuments formata os documentos no padrão esperado pelo LLM
func FormatDocuments(docs []schema.Document) string {
	var contextBuilder strings.Builder
	for _, res := range docs {
		contextBuilder.WriteString(res.PageContent)
		contextBuilder.WriteString("\n\n---\n\n")
	}
	return contextBuilder.String()
}
