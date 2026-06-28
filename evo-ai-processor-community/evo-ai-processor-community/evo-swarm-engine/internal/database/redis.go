package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis inicializa a conexão com o Redis.
func ConnectRedis() (*redis.Client, error) {
	// Puxando configurações via ambiente com default para dev local
	pass := os.Getenv("REDIS_PASSWORD")
	if pass == "" {
		pass = "evoai_redis_pass"
	}

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0, // DB padrão
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao Redis (Memory Engine): %w", err)
	}

	log.Println("✅ Conectado ao Redis com sucesso! (Memory Engine Ativo)")
	return rdb, nil
}
