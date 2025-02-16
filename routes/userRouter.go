package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// ! UserRoutes registers user-related routes
func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/signup", controller.SignUp)
		userGroup.POST("/login", controller.Login)
		userGroup.GET("/", controller.GetUsers)
		userGroup.GET("/:user_id", controller.GetUserByID)
	}
}
