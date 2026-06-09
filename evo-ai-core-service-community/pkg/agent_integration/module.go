package agent_integration

import (
	"evo-ai-core-service/pkg/agent_integration/handler"
	"evo-ai-core-service/pkg/agent_integration/repository"
	"evo-ai-core-service/pkg/agent_integration/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(db *gorm.DB, router gin.IRouter) {
	// Initialize repository
	agentIntegrationRepository := repository.NewAgentIntegrationRepository(db)

	// Initialize service
	agentIntegrationService := service.NewAgentIntegrationService(agentIntegrationRepository)

	// Initialize handler and register routes
	agentIntegrationHandler := handler.NewAgentIntegrationHandler(agentIntegrationService)
	agentIntegrationHandler.RegisterRoutesMiddleware(router)
}
