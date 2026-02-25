package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware validates JWT tokens for protected routes
func Middleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing authorization header", "data": nil})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid authorization format", "data": nil})
			return
		}

		claims, err := ValidateToken(parts[1], secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid or expired token", "data": nil})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
