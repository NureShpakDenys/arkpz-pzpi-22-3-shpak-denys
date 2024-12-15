package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	log *slog.Logger,
	authService services.AuthService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Bearer token missing"})
			c.Abort()
			return
		}

		_, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Set("token", tokenString)
		c.Next()
	}
}
