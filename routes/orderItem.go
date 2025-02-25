package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// ! OrderItemRoutes registers order item-related routes
func OrderItemRoutes(router *gin.Engine) {
	orderItemGroup := router.Group("/orderItems")
	{
		orderItemGroup.GET("/", controller.GetOrderItems())                  //? Get all order items
		orderItemGroup.GET("/:orderItem_id", controller.GetOrderItem())      //? Get a specific order item by ID
		orderItemGroup.POST("/", controller.CreateOrderItem())               //? Create a new order item
		orderItemGroup.PATCH("/:orderItem_id", controller.UpdateOrderItem()) //? Update an existing order item
		// orderItemGroup.DELETE("/:orderItem_id", controller.DeleteOrderItem) //? Delete an order item
	}

	//! This route is registered separately at the root level
	router.GET("/orderItem-order/:order_id", controller.GetOrderItemsByOrder()) //? Get order items by order ID
}
