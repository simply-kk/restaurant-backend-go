package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// ! TableRoutes registers table-related routes
func TableRoutes(router *gin.Engine) {
	tableGroup := router.Group("/tables")
	{
		tableGroup.GET("/", controller.Gettables)               //? Get all tables
		tableGroup.GET("/:table_id", controller.Gettable)       //? Get table by ID
		tableGroup.POST("/", controller.Createtable)            //? Create a new table
		tableGroup.PUT("/:table_id", controller.Updatetable)    //? Update an existing table
		tableGroup.DELETE("/:table_id", controller.Deletetable) //? Delete an table
	}
}
