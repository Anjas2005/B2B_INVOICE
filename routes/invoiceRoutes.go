package invoiceRoutes

import (
	invoiceHandlers "github.com/Anjas2005/bookmarks-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterInvoice(router *gin.Engine, pool *pgxpool.Pool) {
	invoiceHandler := invoiceHandlers.InvoiceHandler{
		DB: pool,
	}

	router.POST("/invoices", invoiceHandler.CreateInvoice)
}
