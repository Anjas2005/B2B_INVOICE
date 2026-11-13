package InvoiceHandlers

import (
	"errors"
	"net/http"

	InvoiceModels "github.com/Anjas2005/bookmarks-api/models"
	"github.com/Anjas2005/bookmarks-api/services"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	Service *services.InvoiceService
}

// CreateInvoice handles batch invoice creation.
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {

	// Receive and decode JSON
	var requests []InvoiceModels.Invoice

	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON request",
			"details": err.Error(),
		})
		return
	}

	// Send invoices to the service
	invoices, err := h.Service.CreateInvoices(
		c.Request.Context(),
		requests,
	)

	// Handle error returned by services
	if err != nil {
		// Check wheather this is our custom validation error.
		var validationErr *services.InvoiceValidationError

		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, validationErr)
			return
		}
		// Any other error is treated as an internal server error.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create invoices",
		})
		return
	}

	// Return successful response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Invoices created successfully",
		"count":   len(invoices),
	})

}
