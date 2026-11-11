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

// TODO Make this such that the api can accept multiple invoices and store them
// and also give the due invoices as a result

// func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
// 	// Creating an empty invoice list
// 	var invoices []InvoiceModels.Invoice

// 	// Reading the incoming JSON
// 	if err := c.ShouldBindJSON(&invoices); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Length needs to be minimum 1
// 	if len(invoices) == 0{
// 		c.JSON(http.StatusBadRequest,gin.H{
// 			"error":"At least one invoice is required",
// 		})
// 		return
// 	}

// 	// Create a var for storing upda

// 	for i := range invoices {

// 		// Check whether the customer has overpaid
// 		if invoice.PaidAmount > invoice.RequiredAmount {
// 			c.JSON(400, gin.H{
// 				"error": "Paid amount cannot exceed required amount",
// 			})
// 			return
// 		}

// 		// Calculate transaction status
// 		if invoice.PaidAmount == invoice.RequiredAmount {

// 			invoice.TransactionStatus = InvoiceModels.Fulfilled

// 			completedAt := time.Now().UTC()
// 			invoice.TransactionCompletionDate = &completedAt

// 		} else {

// 			invoice.TransactionStatus = InvoiceModels.Pending
// 			invoice.TransactionCompletionDate = nil
// 		}

// 		// Parse due date
// 		// TODO: Must be able to parse and undersand all date time formats
// 		dueDate, err := time.Parse("2006-01-02", invoice.DueDate)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date"})
// 			return
// 		}
// 		// if no errors then invoices duedate is made the correct date
// 		invoice.DueDate = dueDate
// 	}

// 	query := `
// 	INSERT INTO invoices (
// 	customer,
// 	invoice_no,
// 	required_amount,
// 	paid_amount,
// 	due_date,
// 	transaction_status,
// 	transaction_completion_date
// 	)
// 	VALUES ($1,$2,$3,$4,$5,$6,$7)
// 	RETURNING id
// 	`

// 	// Inesert the Invoice into DB
// 	err = h.DB.Query(
// 		c.Request.Context(),
// 		query,
// 		invoice.Customer,
// 		invoice.InvoiceNo,
// 		invoice.RequiredAmount,
// 		invoice.PaidAmount,
// 		dueDate,
// 		invoice.TransactionStatus,
// 		invoice.TransactionCompletionDate,
// 	).Scan(&invoice.ID)

// 	if err != nil {
// 		log.Printf("Failed to create invoice: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to create invoice",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, invoice)
// }

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
		// // Convert the validated date string into time.Time for PostgreSQL
		// dueDate, err := time.Parse("2006-01-02", invoices[i].DueDate)
		// if err != nil {
		// 	c.JSON(400, gin.H{"error": "Invalid due date"})
		// 	return
		// }

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
