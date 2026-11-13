package services

import InvoiceModels "github.com/Anjas2005/bookmarks-api/models"

type InvoiceValidationError struct {
	Issues []InvoiceModels.InvoiceIssue `json:"issues"`
}

func (e *InvoiceValidationError) Error() string {
	return "invoice validation failed"
}
