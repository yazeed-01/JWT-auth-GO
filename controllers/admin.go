package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminDashboard(c *gin.Context) {
	// Handle the admin dashboard or other admin-specific logic
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the admin dashboard"})
	// all work on -> middleware/User_Role.go
}
