package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisClient é a instância global do Redis.
var RedisClient *redis.Client

// Ctx é o contexto padrão para operações do Redis.
var Ctx = context.Background()

// ConnectRedis inicializa a conexão com o Redis.
func ConnectRedis() {
	// Puxando configurações via ambiente com default para dev local
	pass := os.Getenv("REDIS_PASSWORD")
	if pass == "" {
		pass = "evoai_redis_pass"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: pass,
		DB:       0, // DB padrão
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Falha crítica: Não foi possível conectar ao Redis (Memory Engine): %v", err)
	}

	log.Println("✅ Conectado ao Redis com sucesso! (Memory Engine Ativo)")
	RedisClient = rdb
}
