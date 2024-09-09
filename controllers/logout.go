package controllers

import (
	"JWTauth/initializers"
	"JWTauth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Get the token from the Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Check if the refresh token exists in the database (my strategy)
	var refreshToken models.RefreshToken
	if err := initializers.DB.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Extract user ID from the refresh token (my strategy)
	userID := refreshToken.UserID

	// Mark existing access tokens as dead (my strategy)
	var existingAccessTokens []models.AccessToken
	if err := initializers.DB.Where("user_id = ?", userID).Find(&existingAccessTokens).Error; err == nil {
		for _, token := range existingAccessTokens {
			token.IsDead = true
			if err := initializers.DB.Save(&token).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update existing access token"})
				return
			}
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find access tokens"})
		return
	}

	// Mark the refresh token as dead (my strategy)
	refreshToken.IsDead = true
	if err := initializers.DB.Save(&refreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update refresh token"})
		return
	}

	// Update the user's IsOnline status to false
	var user models.User
	if err := initializers.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
		return
	}
	user.IsOnline = false
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user status"})
		return
	}
	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
