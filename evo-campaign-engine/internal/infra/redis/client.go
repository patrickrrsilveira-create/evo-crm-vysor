package redis

import (
	"context"
	"evo-campaign-engine/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Printf("redis connected: %s:%s db=%d", cfg.Host, cfg.Port, cfg.DB)
	return client
}
