package httputils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractBearerToken extrai o token Bearer do header Authorization
func ExtractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Verificar formato Bearer
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(authHeader, "Bearer ")
}

