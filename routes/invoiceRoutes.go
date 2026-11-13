package invoiceroutes

import (
	InvoiceHandlers "github.com/Anjas2005/bookmarks-api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterInvoice(router *gin.Engine,
	handler *InvoiceHandlers.InvoiceHandler) {

	router.POST("/invoices", handler.CreateInvoice)
}
