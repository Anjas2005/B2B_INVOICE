package InvoiceModels

import "time"

type TransactionStatus string

const (
	Pending   TransactionStatus = "pending"
	Fulfilled TransactionStatus = "fulfilled"
)

type Invoice struct {
	ID                        int64             `json:"id"`
	Customer                  string            `json:"customer" binding:"required,max=50"`
	InvoiceNo                 string            `json:"invoice_no" binding:"required"`
	RequiredAmount            int64             `json:"required_amount" binding:"required,gte=0"`
	PaidAmount                int64             `json:"paid_amount" binding:"gte=0"`
	DueDate                   string            `json:"due_date" binding:"required,datetime=2006-01-02"`
	TransactionStatus         TransactionStatus `json:"transaction_status"`
	TransactionCompletionDate *time.Time        `json:"transaction_completion_date"`
}
