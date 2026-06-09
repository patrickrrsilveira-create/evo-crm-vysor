package service

import (
	"context"
	"evo-ai-core-service/internal/config"
	"evo-ai-core-service/pkg/core/model"
	"time"

	"gorm.io/gorm"
)

// HealthService interface defines health check operations
type HealthService interface {
	CheckHealth(ctx context.Context) *model.OverallHealth
	CheckReadiness(ctx context.Context) *model.OverallHealth
}

// healthService implements the HealthService interface
type healthService struct {
	healthChecker HealthChecker
}

// NewHealthService creates a new health service
func NewHealthService(db *gorm.DB, config *config.HealthCheckConfig) HealthService {
	healthChecker := NewHealthChecker(db, config)

	return &healthService{
		healthChecker: healthChecker,
	}
}

// CheckHealth performs comprehensive health checks
func (s *healthService) CheckHealth(ctx context.Context) *model.OverallHealth {
	start := time.Now()

	components := []*model.ComponentHealth{
		s.healthChecker.CheckDatabase(ctx),
		s.healthChecker.CheckDatabaseConnections(ctx),
		s.healthChecker.CheckMemory(ctx),
	}

	return model.NewOverallHealth(components, time.Since(start))
}

// CheckReadiness performs readiness checks (lighter than health checks)
func (s *healthService) CheckReadiness(ctx context.Context) *model.OverallHealth {
	start := time.Now()

	components := []*model.ComponentHealth{
		s.healthChecker.CheckDatabasePing(ctx),
	}

	return model.NewOverallHealth(components, time.Since(start))
}
