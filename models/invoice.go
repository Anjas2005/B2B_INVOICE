package InvoiceModels

// The data model of the invoice is as follows
// "customer": "XYZ Stores",
// "invoice_no": "INV-102",
// "amount": 8000,
// "due_date": "2026-11-20"

type Invoice struct {
	ID         int64  `json:"id"`
	Customer   string `json:"customer" binding:"required,max=50"`
	Invoice_No string `json:"invoice_no" binding:"required"`
	Amount     int    `json:"amount" binding:"required,gt=0"`
	Due_Date   string `json:"due_date" binding:"required,datetime=2006-01-02"`
}
