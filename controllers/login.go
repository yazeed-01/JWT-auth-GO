package controllers

import (
	"JWTauth/initializers"
	"JWTauth/models"
	"JWTauth/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	// Get the user data from the JSON body
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if email exists and password(hash) is correct
	var user models.User
	if err := initializers.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}

	// Mark existing access tokens as dead (my strategy)
	var existingAccessTokens []models.AccessToken
	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&existingAccessTokens).Error; err == nil {
		for _, token := range existingAccessTokens {
			token.IsDead = true
			if err := initializers.DB.Save(&token).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update existing access token"})
				return
			}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching existing access tokens"})
		return
	}

	// Mark existing refresh tokens as dead (my strategy)
	var existingRefreshTokens []models.RefreshToken
	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&existingRefreshTokens).Error; err == nil {
		for _, token := range existingRefreshTokens {
			token.IsDead = true
			if err := initializers.DB.Save(&token).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update existing refresh token"})
				return
			}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching existing refresh tokens"})
		return
	}

	// Optionally delete tokens (if needed) (my strategy)
	// if err := initializers.DB.Where("user_id = ?", user.ID).Delete(&models.AccessToken{}).Error; err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete access tokens"})
	//     return
	// }
	// if err := initializers.DB.Where("user_id = ?", user.ID).Delete(&models.RefreshToken{}).Error; err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete refresh tokens"})
	//     return
	// }

	// Insert the new access token into the database (my strategy)
	if err := initializers.DB.Create(&models.AccessToken{UserID: user.ID, Token: accessToken}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create new access token"})
		return
	}

	// Insert the new refresh token into the database (my strategy)
	if err := initializers.DB.Create(&models.RefreshToken{UserID: user.ID, Token: refreshToken}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create new refresh token"})
		return
	}

	// Make IsOnline true for the user
	user.IsOnline = true
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user status"})
		return
	}

	// Return the access token and refresh token
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
