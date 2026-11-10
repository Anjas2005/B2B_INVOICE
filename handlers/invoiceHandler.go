package invoiceHandlers

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
	// Creating an empty invoice
	var invoice InvoiceModels.Invoice

	// Reading the incoming JSON
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Check whether the customer has overpaid
	if invoice.PaidAmount > invoice.RequiredAmount {
		c.JSON(400, gin.H{
			"error": "Paid amount cannot exceed required amount",
		})
		return
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

	// Parse due date
	// TODO: Must be able to parse and undersand all date time formats
	dueDate, err := time.Parse("2006-01-02", invoice.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date"})
		return
	}

	// type Invoice struct {
	// 	ID                        int64             `json:"id"`
	// 	Customer                  string            `json:"customer" binding:"required,max=50"`
	// 	InvoiceNo                 string            `json:"invoice_no" binding:"required"`
	// 	RequiredAmount            int64             `json:"required_amount" binding:"required,gte=0"`
	// 	PaidAmount                int64             `json:"paid_amount" binding:"gte=0"`
	// 	DueDate                   string            `json:"due_date" binding:"required,datetime=2006-01-02"`
	// 	TransactionStatus         TransactionStatus `json:"transaction_status"`
	// 	TransactionCompletionDate *time.Time        `json:"transaction_completion_date"`
	// }
	// Query to insert the invocie into db
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
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`

	// Inesert the Invoice into DB
	err = h.DB.QueryRow(
		c.Request.Context(),
		query,
		invoice.Customer,
		invoice.InvoiceNo,
		invoice.RequiredAmount,
		invoice.PaidAmount,
		dueDate,
		invoice.TransactionStatus,
		invoice.TransactionCompletionDate,
	).Scan(&invoice.ID)

	if err != nil {
		log.Printf("Failed to create invoice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create invoice",
		})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}
