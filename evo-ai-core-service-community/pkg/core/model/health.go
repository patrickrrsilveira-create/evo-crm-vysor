package model

import "time"

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// ComponentHealth represents the health of a single component
type ComponentHealth struct {
	Name      string                 `json:"name"`
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  string                 `json:"duration"` // Changed to string for Swagger compatibility
}

// OverallHealth represents the overall system health
type OverallHealth struct {
	Status     HealthStatus       `json:"status"`
	Timestamp  time.Time          `json:"timestamp"`
	Duration   string             `json:"duration"` // Changed to string for Swagger compatibility
	Components []*ComponentHealth `json:"components"`
}

// NewComponentHealth creates a new ComponentHealth with basic info
func NewComponentHealth(name string) *ComponentHealth {
	return &ComponentHealth{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Details:   make(map[string]interface{}),
	}
}

// SetHealthy marks the component as healthy
func (c *ComponentHealth) SetHealthy(message string, duration time.Duration) {
	c.Status = HealthStatusHealthy
	c.Message = message
	c.Duration = duration.String()
}

// SetUnhealthy marks the component as unhealthy
func (c *ComponentHealth) SetUnhealthy(message string, duration time.Duration) {
	c.Status = HealthStatusUnhealthy
	c.Message = message
	c.Duration = duration.String()
}

// SetDegraded marks the component as degraded
func (c *ComponentHealth) SetDegraded(message string, duration time.Duration) {
	c.Status = HealthStatusDegraded
	c.Message = message
	c.Duration = duration.String()
}

// AddDetail adds a detail to the component health
func (c *ComponentHealth) AddDetail(key string, value interface{}) {
	c.Details[key] = value
}

// NewOverallHealth creates a new OverallHealth
func NewOverallHealth(components []*ComponentHealth, totalDuration time.Duration) *OverallHealth {
	overallStatus := determineOverallStatus(components)

	return &OverallHealth{
		Status:     overallStatus,
		Timestamp:  time.Now().UTC(),
		Duration:   totalDuration.String(),
		Components: components,
	}
}

// determineOverallStatus calculates the overall status based on components
func determineOverallStatus(components []*ComponentHealth) HealthStatus {
	overallStatus := HealthStatusHealthy

	for _, component := range components {
		if component.Status == HealthStatusUnhealthy {
			return HealthStatusUnhealthy
		}

		if component.Status == HealthStatusDegraded && overallStatus != HealthStatusUnhealthy {
			overallStatus = HealthStatusDegraded
		}
	}

	return overallStatus
}
