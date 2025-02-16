package routes

import (
	"github.com/gin-gonic/gin"
	controllers "golang-restaurant-management/controllers"
)

// ! InvoiceRoutes registers invoice-related routes
func InvoiceRoutes(router *gin.Engine) {
	invoiceGroup := router.Group("/invoices")
	{
		invoiceGroup.GET("/", controllers.GetInvoices)                 //? Get all invoices
		invoiceGroup.GET("/:invoice_id", controllers.GetInvoice)       //? Get invoice by ID
		invoiceGroup.POST("/", controllers.CreateInvoice)              //? Create a new invoice
		invoiceGroup.PUT("/:invoice_id", controllers.UpdateInvoice)    //? Update an invoice
		invoiceGroup.DELETE("/:invoice_id", controllers.DeleteInvoice) //? Delete an invoice
	}
}
