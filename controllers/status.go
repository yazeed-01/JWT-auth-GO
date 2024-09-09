package controllers

import (
	"JWTauth/initializers"
	"JWTauth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckStatus checks if the user is logged in or not
func CheckStatus(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	token := authHeader
	// parse the access token
	/* claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		return
	}
	*/

	// Check if the access token exists in the database (my strategy)
	var accessToken models.AccessToken
	if err := initializers.DB.Where("token = ?", token).First(&accessToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		return
	}

	// Check if token IsDead (my strategy)
	if accessToken.IsDead {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
		return
	}

	// Use the ID from the access token to find the user
	var user models.User
	if err := initializers.DB.Where("id = ?", accessToken.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user's online status
	c.JSON(http.StatusOK, gin.H{"is_online": user.IsOnline})
}
