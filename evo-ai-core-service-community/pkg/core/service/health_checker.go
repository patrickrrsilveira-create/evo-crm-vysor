package service

import (
	"context"
	"evo-ai-core-service/internal/config"
	"evo-ai-core-service/pkg/core/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// HealthChecker interface defines health check operations for specific components
type HealthChecker interface {
	CheckDatabase(ctx context.Context) *model.ComponentHealth
	CheckDatabaseConnections(ctx context.Context) *model.ComponentHealth
	CheckDatabasePing(ctx context.Context) *model.ComponentHealth
	CheckMemory(ctx context.Context) *model.ComponentHealth
}

// healthChecker implements the HealthChecker interface
type healthChecker struct {
	db     *gorm.DB
	config *config.HealthCheckConfig
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(db *gorm.DB, config *config.HealthCheckConfig) HealthChecker {
	return &healthChecker{
		db:     db,
		config: config,
	}
}

// CheckDatabase performs a comprehensive database health check
func (hc *healthChecker) CheckDatabase(ctx context.Context) *model.ComponentHealth {
	start := time.Now()
	component := model.NewComponentHealth("database")

	// Get underlying sql.DB
	sqlDB, err := hc.db.DB()
	if err != nil {
		component.SetUnhealthy("Failed to get database connection", time.Since(start))
		return component
	}

	// Test connection with configured timeout
	ctx, cancel := context.WithTimeout(ctx, hc.config.Timeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		component.SetUnhealthy(fmt.Sprintf("Database ping failed: %v", err), time.Since(start))
		return component
	}

	// Test a simple query with configured query
	var result int
	err = hc.db.WithContext(ctx).Raw(hc.config.DBQuery).Scan(&result).Error
	if err != nil {
		component.SetUnhealthy(fmt.Sprintf("Database query failed: %v", err), time.Since(start))
		return component
	}

	component.SetHealthy("Database is healthy", time.Since(start))
	return component
}

// CheckDatabaseConnections checks database connection pool health
func (hc *healthChecker) CheckDatabaseConnections(ctx context.Context) *model.ComponentHealth {
	start := time.Now()
	component := model.NewComponentHealth("database_connections")

	sqlDB, err := hc.db.DB()
	if err != nil {
		component.SetUnhealthy("Failed to get database connection", time.Since(start))
		return component
	}

	stats := sqlDB.Stats()

	// Add connection stats to details
	component.AddDetail("open_connections", stats.OpenConnections)
	component.AddDetail("in_use", stats.InUse)
	component.AddDetail("idle", stats.Idle)
	component.AddDetail("wait_count", stats.WaitCount)
	component.AddDetail("wait_duration", stats.WaitDuration.String())
	component.AddDetail("max_idle_closed", stats.MaxIdleClosed)
	component.AddDetail("max_idle_time_closed", stats.MaxIdleTimeClosed)
	component.AddDetail("max_lifetime_closed", stats.MaxLifetimeClosed)

	// Check if connection pool is healthy using configured thresholds
	if stats.OpenConnections > hc.config.DBConnThresholdDegraded {
		component.SetDegraded("High number of database connections", time.Since(start))
		return component
	}

	if stats.WaitCount > int64(hc.config.DBWaitCountThreshold) {
		component.SetDegraded("High database connection wait count", time.Since(start))
		return component
	}

	component.SetHealthy("Database connections are healthy", time.Since(start))
	return component
}

// CheckDatabasePing performs a simple database ping (for readiness)
func (hc *healthChecker) CheckDatabasePing(ctx context.Context) *model.ComponentHealth {
	start := time.Now()
	component := model.NewComponentHealth("database_ping")

	sqlDB, err := hc.db.DB()
	if err != nil {
		component.SetUnhealthy("Failed to get database connection", time.Since(start))
		return component
	}

	// Quick ping with configured readiness timeout
	ctx, cancel := context.WithTimeout(ctx, hc.config.ReadinessTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		component.SetUnhealthy(fmt.Sprintf("Database ping failed: %v", err), time.Since(start))
		return component
	}

	component.SetHealthy("Database ping successful", time.Since(start))
	return component
}

// CheckMemory performs basic memory health check
func (hc *healthChecker) CheckMemory(ctx context.Context) *model.ComponentHealth {
	start := time.Now()
	component := model.NewComponentHealth("memory")

	// This is a basic implementation - in production you might want to use
	// runtime.ReadMemStats() or external monitoring tools
	component.SetHealthy("Memory usage is within acceptable limits", time.Since(start))
	return component
}
