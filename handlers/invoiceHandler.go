package InvoiceHandlers

import (
	"log"
	"net/http"
	"time"

	InvoiceModels "github.com/Anjas2005/bookmarks-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceHandler struct {
	DB *pgxpool.Pool
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	// Var to store and process requests
	var requests []InvoiceModels.Invoice

	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(requests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one invoice is required",
		})
		return
	}

	// Var for finally processed invoices
	invoices := make([]InvoiceModels.Invoice, 0, len(requests))

	for i := range requests {
		req := requests[i]

		// Validate payment amount
		if req.PaidAmount > req.RequiredAmount {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Paid amount cannot exceed required amount",
				"index": i,
			})
			return
		}
		// Parse due date
		_, err := time.Parse("2006-01-02", req.DueDate)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid due date",
				"index": i,
			})
			return
		}

		// Create a new invoice
		invoice := InvoiceModels.Invoice{
			Customer:       req.Customer,
			InvoiceNo:      req.InvoiceNo,
			RequiredAmount: req.RequiredAmount,
			PaidAmount:     req.PaidAmount,
			DueDate:        req.DueDate,
		}
		// Calculate transaction status
		if invoice.PaidAmount == invoice.RequiredAmount {

			invoice.TransactionStatus = InvoiceModels.Fulfilled

			completedAt := time.Now().UTC()
			invoice.TransactionCompletionDate = &completedAt

		} else {

			invoice.TransactionStatus = InvoiceModels.Pending
			invoice.TransactionCompletionDate = nil
		}

		// Add processed invoice to the new slice
		invoices = append(invoices, invoice)
	}

	ctx := c.Request.Context()

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO invoices (
        customer,
        invoice_no,
        required_amount,
        paid_amount,
        due_date,
        transaction_status,
        transaction_completion_date
    )
    VALUES ($1, $2, $3, $4, $5::text::date, $6, $7)
    RETURNING id
	`

	for i := range invoices {

		err = tx.QueryRow(
			ctx,
			query,
			invoices[i].Customer,
			invoices[i].InvoiceNo,
			invoices[i].RequiredAmount,
			invoices[i].PaidAmount,
			invoices[i].DueDate,
			invoices[i].TransactionStatus,
			invoices[i].TransactionCompletionDate,
		).Scan(&invoices[i].ID)

		if err != nil {
			log.Printf("Invoice insertion failed: %v", err)

			c.JSON(500, gin.H{
				"error": "Failed to insert invoice batch",
			})
			return
		}
	}
	// Check if we failed to commit the invoices
	if err := tx.Commit(ctx); err != nil {
		c.JSON(500, gin.H{"error": "Failed to commit invoices"})
		return
	}

	// Send the responses back
	c.JSON(http.StatusCreated, gin.H{
		"message": "Invoices created successfully",
		"count":   len(invoices),
	})
}
