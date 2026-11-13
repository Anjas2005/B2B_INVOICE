package InvoiceModels

type InvoiceIssue struct {
	Index   int     `json:"index"`
	Field   string  `json:"field"`
	Message string  `json:"message"`
	Invoice Invoice `json:"invoice"`
}
