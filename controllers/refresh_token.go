package controllers

import (
	"JWTauth/initializers"
	"JWTauth/models"
	"JWTauth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Parse the refresh token
	claims, err := utils.ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Check if the refresh token exists in the database (my strategy)
	var refreshToken models.RefreshToken
	if err := initializers.DB.First(&refreshToken, "user_id = ?", claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Check if the refresh token is marked as dead (my strategy)
	if refreshToken.IsDead {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Verify if the user exists (my strategy)
	var user models.User
	if err := initializers.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if the user has an access token in the database (my strategy)
	var accessToken models.AccessToken
	if err := initializers.DB.First(&accessToken, "user_id = ?", claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		return
	}

	// Check if the access token is marked as dead (my strategy)
	if accessToken.IsDead {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		return
	}

	// Generate a new access token
	newAccessToken, _, err := utils.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	// Delete the old access token from the database (my strategy)
	if err := initializers.DB.Delete(&accessToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete old access token"})
		return
	}

	// Insert the new access token into the database (my strategy)
	if err := initializers.DB.Create(&models.AccessToken{UserID: user.ID, Token: newAccessToken}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save new access token"})
		return
	}

	// Return the new access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
