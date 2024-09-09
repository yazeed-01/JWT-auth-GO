package routes

import (
	"JWTauth/controllers"
	"JWTauth/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.POST("/refresh-token", controllers.RefreshToken)

	r.GET("/user-info", middleware.AuthMiddleware(), controllers.UserInfo)
	r.GET("/check-status", middleware.AuthMiddleware(), controllers.CheckStatus)

	// other way to apply middleware:
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AdminRoleRequired())
	{
		adminRoutes.GET("/", controllers.AdminDashboard)
		// other routes ....
	}

	return r

}
