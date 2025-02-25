package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

//! UserRoutes registers user-related routes
func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/signup", controller.SignUp())   //? Register a new user
		userGroup.POST("/login", controller.Login())     //? Authenticate a user
		userGroup.GET("/", controller.GetUsers())        //? Get all users
		userGroup.GET("/:user_id", controller.GetUser()) //? Get user by ID
	}
}
