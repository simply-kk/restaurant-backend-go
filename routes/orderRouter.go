package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// ! OrderRoutes registers order-related routes
func OrderRoutes(router *gin.Engine) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.GET("/", controller.GetOrders)               //? Get all orders
		orderGroup.GET("/:order_id", controller.GetOrder)       //? Get order by ID
		orderGroup.POST("/", controller.CreateOrder)            //? Create a new order
		orderGroup.PUT("/:order_id", controller.UpdateOrder)    //? Update an existing order
		orderGroup.DELETE("/:order_id", controller.DeleteOrder) //? Delete an order
	}
}
