package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port  string
	DB    DBConfig
	Redis RedisConfig

	EvolutionAPIURL string
	N8NWebhookURL   string
	CRMBaseURL      string
	CRMAPIToken     string

	WorkerCount     int
	SchedulerTick   time.Duration
	MaxRetries      int
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load() *Config {
	return &Config{
		Port: envOrDefault("PORT", "9090"),
		DB: DBConfig{
			Host:            envOrPanic("DB_HOST"),
			Port:            envOrDefault("DB_PORT", "5432"),
			User:            envOrPanic("DB_USER"),
			Password:        envOrPanic("DB_PASSWORD"),
			DBName:          envOrPanic("DB_NAME"),
			SSLMode:         envOrDefault("DB_SSLMODE", "disable"),
			MaxIdleConns:    envOrDefaultInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    envOrDefaultInt("DB_MAX_OPEN_CONNS", 50),
			ConnMaxLifetime: envOrDefaultDuration("DB_CONN_MAX_LIFETIME", time.Hour),
			ConnMaxIdleTime: envOrDefaultDuration("DB_CONN_MAX_IDLE_TIME", 30*time.Minute),
		},
		Redis: RedisConfig{
			Host:     envOrPanic("REDIS_HOST"),
			Port:     envOrDefault("REDIS_PORT", "6379"),
			Password: envOrDefault("REDIS_PASSWORD", ""),
			DB:       envOrDefaultInt("REDIS_DB", 6),
		},

		EvolutionAPIURL: envOrDefault("EVOLUTION_API_URL", ""),
		N8NWebhookURL:   envOrDefault("N8N_WEBHOOK_URL", ""),
		CRMBaseURL:      envOrDefault("CRM_BASE_URL", ""),
		CRMAPIToken:     envOrDefault("CRM_API_TOKEN", ""),

		WorkerCount:     envOrDefaultInt("WORKER_COUNT", 4),
		SchedulerTick:   envOrDefaultDuration("SCHEDULER_TICK", 2*time.Second),
		MaxRetries:      envOrDefaultInt("MAX_RETRIES", 3),
		ShutdownTimeout: envOrDefaultDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
	}
}

func envOrPanic(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return strings.TrimSpace(v)
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return strings.TrimSpace(v)
}

func envOrDefaultInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		log.Printf("warning: invalid int for %s=%q, using %d", key, v, fallback)
		return fallback
	}
	return n
}

func envOrDefaultDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		log.Printf("warning: invalid duration for %s=%q, using %s", key, v, fallback)
		return fallback
	}
	return d
}
