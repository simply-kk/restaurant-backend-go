package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// ! FoodRoutes registers food-related routes
func FoodRoutes(router *gin.Engine) {
	foodGroup := router.Group("/food")
	{
		foodGroup.GET("/", controller.GetFoods)              //? Get all foods
		foodGroup.GET("/:food_id", controller.GetFood)       //? Get food by ID
		foodGroup.POST("/", controller.CreateFood)           //? Create a new food item
		foodGroup.PATCH("/:food_id", controller.UpdateFood)  //? Update an existing food item
		foodGroup.DELETE("/:food_id", controller.DeleteFood) //? Delete food item
	}
}
