package controllers

import (
	"JWTauth/initializers"
	"JWTauth/models"
	"JWTauth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	token := authHeader

	// Parse the access token
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		return
	}

	// Check if the access token exists in the database (my strategy)
	var accessToken models.AccessToken
	if err := initializers.DB.Where("user_id = ?", claims.UserID).First(&accessToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		return
	}

	// Check if the access token is marked as dead (my strategy)
	if accessToken.IsDead {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
		return
	}

	// Retrieve user information
	var user models.User
	if err := initializers.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Respond with user information
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}
