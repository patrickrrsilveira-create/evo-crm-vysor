package memory
	
	import (
		"context"
		"encoding/json"
		"fmt"
		"time"
	
		"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
		"github.com/redis/go-redis/v9"
	)
	
	// ShortMemoryManager lida com a memória de curto prazo, baseada no conceito de cache hiper-rápido do Redis
	type ShortMemoryManager struct {
		Client *redis.Client
		TTL    time.Duration
	}
	
// NewShortMemoryManager cria um novo gerenciador de memória usando o client Redis injetado.
func NewShortMemoryManager(client *redis.Client) *ShortMemoryManager {
	return &ShortMemoryManager{
		Client: client,
		TTL:    24 * time.Hour, // O short memory expira rápido
	}
}
	
	// AddMessageToHistory insere uma nova mensagem na memória de curto prazo associada a uma conversa
	func (r *ShortMemoryManager) AddMessageToHistory(ctx context.Context, conversationID string, msg models.LLMMessage) error {
		key := fmt.Sprintf("short_memory:conversation:%s", conversationID)
		
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
	
		// LPush adiciona a mensagem no início da lista (O(1))
		pipe := r.Client.Pipeline()
		pipe.LPush(ctx, key, msgBytes)
		// LTrim mantém apenas as últimas 20 mensagens (janela de curto prazo)
		pipe.LTrim(ctx, key, 0, 19)
		// Renova o TTL
		pipe.Expire(ctx, key, r.TTL)
		
		_, err = pipe.Exec(ctx)
		return err
	}
	
	// GetRecentHistory busca as mensagens mais recentes (Short Memory)
	func (r *ShortMemoryManager) GetRecentHistory(ctx context.Context, conversationID string) ([]models.LLMMessage, error) {
		key := fmt.Sprintf("short_memory:conversation:%s", conversationID)
	
		// Redis LRange retorna a partir do index 0 até o final (-1).
		// Como fizemos LPush, as mais novas estão no index 0. Precisaremos inverter a ordem para injetar na LLM.
		result, err := r.Client.LRange(ctx, key, 0, -1).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
	
		if len(result) == 0 {
			return []models.LLMMessage{}, nil
		}
	
		var history []models.LLMMessage
		// Inverte para ficar cronológico (Antigas -> Novas)
		for i := len(result) - 1; i >= 0; i-- {
			var msg models.LLMMessage
			if err := json.Unmarshal([]byte(result[i]), &msg); err == nil {
				history = append(history, msg)
			}
		}
	
		return history, nil
	}
	
	// ClearShortMemory limpa explicitamente a memória de curto prazo (ex: se o usuário pedir para reiniciar o assunto)
	func (r *ShortMemoryManager) ClearShortMemory(ctx context.Context, conversationID string) error {
		key := fmt.Sprintf("short_memory:conversation:%s", conversationID)
		return r.Client.Del(ctx, key).Err()
	}
