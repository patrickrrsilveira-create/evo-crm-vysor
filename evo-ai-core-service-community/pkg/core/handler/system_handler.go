package handler

import (
	"evo-ai-core-service/internal/httpclient/response"
	"evo-ai-core-service/pkg/core/model"
	"evo-ai-core-service/pkg/core/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SystemHandler interface defines the contract for system handlers
type SystemHandler interface {
	RegisterRoutes(router gin.IRouter)
}

// systemHandler implements the SystemHandler interface
type systemHandler struct {
	healthService service.HealthService
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(healthService service.HealthService) SystemHandler {
	return &systemHandler{
		healthService: healthService,
	}
}

// RegisterRoutes registers the routes for the system handler
func (h *systemHandler) RegisterRoutes(router gin.IRouter) {
	router.GET("/health", h.HealthCheck)
	router.GET("/ready", h.ReadyCheck)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// HealthCheck handles the health check request
func (h *systemHandler) HealthCheck(c *gin.Context) {
	health := h.healthService.CheckHealth(c.Request.Context())

	statusCode := http.StatusOK
	if health.Status == model.HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	if health.Status == model.HealthStatusDegraded {
		statusCode = http.StatusOK // Still return 200 for degraded but include status
	}

	response.SuccessResponse(c, gin.H{
		"status": health.Status,
		"service": "evo-ai-core-service",
		"version": "1.0.0",
		"components": health.Components,
	}, "", statusCode)
}

// ReadyCheck handles the ready check request
func (h *systemHandler) ReadyCheck(c *gin.Context) {
	readiness := h.healthService.CheckReadiness(c.Request.Context())

	statusCode := http.StatusOK
	if readiness.Status == model.HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	response.SuccessResponse(c, gin.H{
		"status": readiness.Status,
		"service": "evo-ai-core-service",
		"version": "1.0.0",
		"components": readiness.Components,
	}, "", statusCode)
}
