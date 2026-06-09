package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAdminMiddleware interface {
	GetUserAdminMiddleware() gin.HandlerFunc
}

type userAdminMiddleware struct {
}

func NewUserAdminMiddleware() UserAdminMiddleware {
	return &userAdminMiddleware{}
}

func (u *userAdminMiddleware) GetUserAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := c.GetBool("is_admin")
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is not an admin"})
			return
		}

		c.Next()
	}
}
