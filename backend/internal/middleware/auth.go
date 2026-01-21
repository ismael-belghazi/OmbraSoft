package middleware

import (
	"strings"

	"github.com/bebeb/ombrasoft-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Token manquant"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Format de token invalide"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Token invalide ou expiré"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
