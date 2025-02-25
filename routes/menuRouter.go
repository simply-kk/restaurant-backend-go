package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-restaurant-management/controllers"
)

// ! MenuRoutes registers menu-related route
func MenuRoutes(router *gin.Engine) {
	menuGroup := router.Group("/menus")
	{
		menuGroup.GET("/", controller.GetMenus())              //? Get all menus
		menuGroup.GET("/:menu_id", controller.GetMenu())       //? Get menu by ID
		menuGroup.POST("/", controller.CreateMenu())           //? Create a new menu
		menuGroup.PATCH("/:menu_id", controller.UpdateMenu())    //? Update an menu
		// menuGroup.DELETE("/:menu_id", controller.DeleteMenu) //? Delete an menu
	}
}
