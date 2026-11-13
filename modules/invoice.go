package modules

import (
	InvoiceHandlers "github.com/Anjas2005/bookmarks-api/handlers"
	invoiceroutes "github.com/Anjas2005/bookmarks-api/routes"
	"github.com/Anjas2005/bookmarks-api/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterInvoice(
	router *gin.Engine,
	dbPool *pgxpool.Pool,
) {
	invoiceService := &services.InvoiceService{
		DB: dbPool,
	}

	invoiceHandler := &InvoiceHandlers.InvoiceHandler{
		Service: invoiceService,
	}

	invoiceroutes.RegisterInvoice(router, invoiceHandler)
}
