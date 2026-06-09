package postgres

import (
	"database/sql"
	"evo-ai-core-service/internal/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/lib/pq"
)

func Connect(DBConfig *config.DBConfig) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DBConfig.Host, DBConfig.Port, DBConfig.User, DBConfig.Password, DBConfig.DBName, DBConfig.SSLMode)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	log.Println("Connected to Postgres database.")

	return conn
}

func ConnectGorm(DBConfig *config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DBConfig.Host, DBConfig.Port, DBConfig.User, DBConfig.Password, DBConfig.DBName, DBConfig.SSLMode)

	// Configure GORM with optimized settings
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:                              true,  // Cache prepared statements
		DisableForeignKeyConstraintWhenMigrating: false, // Keep foreign key constraints
	}

	conn, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalf("error getting underlying sql.DB: %v", err)
	}

	// Parse configuration values
	maxIdleConns := parseIntOrDefault(DBConfig.MaxIdleConns, 10)
	maxOpenConns := parseIntOrDefault(DBConfig.MaxOpenConns, 100)
	connMaxLifetime := parseDurationOrDefault(DBConfig.ConnMaxLifetime, time.Hour)
	connMaxIdleTime := parseDurationOrDefault(DBConfig.ConnMaxIdleTime, 30*time.Minute)

	// Configure connection pool with environment values
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	log.Printf("Connected to Postgres database with connection pool: idle=%d, max=%d, lifetime=%s, idle_time=%s",
		maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime)

	return conn
}

// parseIntOrDefault parses string to int with fallback to default
func parseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Invalid integer value: %s, using default: %d", value, defaultValue)
		return defaultValue
	}
	return intValue
}

// parseDurationOrDefault parses string to time.Duration with fallback to default
func parseDurationOrDefault(value string, defaultValue time.Duration) time.Duration {
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Warning: Invalid duration value: %s, using default: %s", value, defaultValue.String())
		return defaultValue
	}
	return duration
}
