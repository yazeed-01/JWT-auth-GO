package middleware

import (
	"JWTauth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the access token from the Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Log the token parsing for debugging purposes
		c.Set("token", token) // Set token in context

		// Parse and validate the token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the UserID in the context for use in other handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
