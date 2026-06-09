package middleware

import (
	"context"
	"fmt"
	"net/http"

	"evo-ai-core-service/internal/httpclient/response"
	"evo-ai-core-service/internal/services"
	"evo-ai-core-service/internal/utils/contextutils"

	"github.com/gin-gonic/gin"
)

// Singleton global para o middleware de permissão
var globalPermissionMiddleware PermissionMiddleware

// InitializePermissionMiddleware inicializa o middleware global
func InitializePermissionMiddleware(evoAuthBaseURL string) {
	globalPermissionMiddleware = NewPermissionMiddleware(evoAuthBaseURL)
}

// GetGlobalPermissionMiddleware retorna o middleware global
func GetGlobalPermissionMiddleware() PermissionMiddleware {
	if globalPermissionMiddleware == nil {
		panic("Permission middleware not initialized. Call InitializePermissionMiddleware first.")
	}
	return globalPermissionMiddleware
}

// PermissionMiddleware interface para validação de permissões
type PermissionMiddleware interface {
	RequirePermission(resource, action string) gin.HandlerFunc
	CheckPermission(authToken, permissionKey string) (bool, error)
	CheckPermissionWithType(authToken, permissionKey, tokenType string) (bool, error)
}

type permissionMiddleware struct {
	authService services.EvoAuthService
}

// NewPermissionMiddleware cria uma nova instância do middleware de permissões
// Delegando toda lógica para EvoAuthService para consistência
func NewPermissionMiddleware(evoAuthBaseURL string) PermissionMiddleware {
	return &permissionMiddleware{
		authService: services.NewEvoAuthService(evoAuthBaseURL),
	}
}

// RequirePermission cria um middleware que valida se o usuário tem a permissão específica
func (p *permissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Construir chave de permissão
		permissionKey := fmt.Sprintf("%s.%s", resource, action)

		// Get token type from context (set by EvoAuthMiddleware)
		tokenType, err := contextutils.GetTokenType(c.Request.Context())
		authToken, err := contextutils.GetToken(c.Request.Context())
		if authToken == "" {
			response.ErrorResponse(c, "ERR_UNAUTHORIZED", "Token must be provided in header", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		hasPermission, err := p.CheckPermissionWithType(authToken, permissionKey, tokenType)

		if err != nil {
			fmt.Printf("Permission: Error checking permission %s: %v\n", permissionKey, err)
			response.ErrorResponse(c, "ERR_INTERNAL_SERVER", "Unable to validate user permissions", nil, http.StatusInternalServerError)
			c.Abort()
			return
		}

		fmt.Printf("Permission: Has permission %s: %v\n", permissionKey, hasPermission)

		if !hasPermission {
			response.ErrorResponse(c, "ERR_FORBIDDEN", "Insufficient permissions", nil, http.StatusForbidden)
			c.Abort()
			return
		}

		fmt.Printf("Permission: Access granted for permission %s\n", permissionKey)
		c.Next()
	}
}

// CheckPermission delegates to EvoAuthService for unified permission handling
func (p *permissionMiddleware) CheckPermission(authToken, permissionKey string) (bool, error) {
	return p.CheckPermissionWithType(authToken, permissionKey, "bearer")
}

// CheckPermissionWithType delegates to EvoAuthService for unified permission handling with specific token type
func (p *permissionMiddleware) CheckPermissionWithType(authToken, permissionKey, tokenType string) (bool, error) {
	return p.authService.CheckPermission(context.Background(), authToken, permissionKey, tokenType)
}
