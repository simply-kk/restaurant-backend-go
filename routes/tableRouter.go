package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// TableRoutes registers table-related routes
func TableRoutes(router *gin.Engine) {
	tableGroup := router.Group("/tables")
	{
		tableGroup.GET("/", controller.GetTables())               //? Get all tables
		tableGroup.GET("/:table_id", controller.GetTable())       //? Get table by ID
		tableGroup.POST("/", controller.CreateTable())            //? Create a new table
		tableGroup.PATCH("/:table_id", controller.UpdateTable())    //? Update a table
		// tableGroup.DELETE("/:table_id", controller.DeleteTable()) //? Delete a table
	}
}
