package services

import (
	"context"
	"fmt"
	"time"

	InvoiceModels "github.com/Anjas2005/bookmarks-api/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceService struct {
	DB *pgxpool.Pool
}

// Function to create invoices
// accepts the context of the requests and requests(type invoice) returns the InvoiceModels and error
func (s *InvoiceService) CreateInvoices(
	ctx context.Context,
	requests []InvoiceModels.Invoice,
) ([]InvoiceModels.Invoice, error) {

	if len(requests) == 0 {
		return nil, fmt.Errorf("at least one invoice is required")
	}

	// make([]type,length,capacity)
	invoices := make([]InvoiceModels.Invoice, 0, len(requests))
	issues := make([]InvoiceModels.InvoiceIssue, 0)

	// Validation
	for i := range requests {
		req := requests[i]
		valid := true

		// Customer validation
		if req.Customer == "" {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "customer",
				Message: "Customer is required",
				Invoice: req,
			})
			valid = false
		} else if len(req.Customer) > 50 {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "customer",
				Message: "Customer cannot exceed 50 characters",
				Invoice: req,
			})
			valid = false
		}

		// Invoice number validation
		if req.InvoiceNo == "" {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "invoice_no",
				Message: "Invoice number is required",
				Invoice: req,
			})
			valid = false
		} else if len(req.InvoiceNo) > 50 {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "invoice_no",
				Message: "Invoice number cannot exceed 50 characters",
				Invoice: req,
			})
			valid = false
		}

		// Required amount validation
		if req.RequiredAmount <= 0 {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "required_amount",
				Message: "Required amount must be greater than zero",
				Invoice: req,
			})
			valid = false
		}

		// Paid amount validation
		if req.PaidAmount < 0 {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "paid_amount",
				Message: "Paid amount cannot be negative",
				Invoice: req,
			})
			valid = false
		} else if req.PaidAmount > req.RequiredAmount {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "paid_amount",
				Message: "Paid amount exceeds required amount",
				Invoice: req,
			})
			valid = false
		}

		// Due date validation
		if _, err := time.Parse("2006-01-02", req.DueDate); err != nil {
			issues = append(issues, InvoiceModels.InvoiceIssue{
				Index:   i,
				Field:   "due_date",
				Message: "Invalid date. Expected YYYY-MM-DD",
				Invoice: req,
			})
			valid = false
		}

		// Don't process an invalid invoice.
		if !valid {
			continue
		}

		// Calculate transaction status.
		invoice := req

		if invoice.PaidAmount == invoice.RequiredAmount {
			invoice.TransactionStatus = InvoiceModels.Fulfilled

			completedAt := time.Now().UTC()
			invoice.TransactionCompletionDate = &completedAt
		} else {
			invoice.TransactionStatus = InvoiceModels.Pending
			invoice.TransactionCompletionDate = nil
		}

		invoices = append(invoices, invoice)
	}

	// Send all errors together.
	if len(issues) > 0 {
		return nil, &InvoiceValidationError{
			Issues: issues,
		}
	}

	// Prepare rows for PostgreSQL COPY.
	rows := make([][]any, 0, len(invoices))

	for _, invoice := range invoices {
		dueDate, _ := time.Parse("2006-01-02", invoice.DueDate)

		rows = append(rows, []any{
			invoice.Customer,
			invoice.InvoiceNo,
			invoice.RequiredAmount,
			invoice.PaidAmount,
			dueDate,
			string(invoice.TransactionStatus),
			invoice.TransactionCompletionDate,
		})
	}
	// Start database transaction.
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin invoice transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert in Bulk using PostgresSQL COPY
	columns := []string{
		"customer",
		"invoice_no",
		"required_amount",
		"paid_amount",
		"due_date",
		"transaction_status",
		"transaction_completion_date",
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"invoices"},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return nil, fmt.Errorf("copy invoices: %w", err)
	}

	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit invoice transaction: %w", err)
	}
	return invoices, nil
}
