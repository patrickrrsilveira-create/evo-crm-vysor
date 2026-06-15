package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisClient é a instância global do Redis.
var RedisClient *redis.Client

// Ctx é o contexto padrão para operações do Redis.
var Ctx = context.Background()

// ConnectRedis inicializa a conexão com o Redis.
func ConnectRedis() {
	// A URL do redis local no docker compose costuma ser 6379 sem senha por padrão neste setup.
	// Podemos puxar de variável de ambiente, mas como é um ambiente de dev com localhost, usamos o padrão.
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "evoai_redis_pass",
		DB:       0,  // DB padrão
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Falha crítica: Não foi possível conectar ao Redis (Memory Engine): %v", err)
	}

	log.Println("✅ Conectado ao Redis com sucesso! (Memory Engine Ativo)")
	RedisClient = rdb
}
