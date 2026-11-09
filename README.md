# Invoice Overdue API

A lightweight REST API for B2B SaaS applications and personal financial analysis.

**Purpose:** Identify overdue customer invoices and calculate how many days they are past their due date.

## Example

**Endpoint:** `POST /invoices/analyze`

**Input:**
```json
{
  "invoices": [
    {
      "customer": "ABC Traders",
      "invoice_no": "INV-101",
      "amount": 15000,
      "due_date": "2026-11-01"
    },
    {
      "customer": "XYZ Stores",
      "invoice_no": "INV-102",
      "amount": 8000,
      "due_date": "2026-11-20"
    }
  ]
}
```

**Output:**
```json
{
  "overdue_count": 1,
  "overdue_amount": 15000,
  "invoices": [
    {
      "customer": "ABC Traders",
      "invoice_no": "INV-101",
      "amount": 15000,
      "days_overdue": 8
    }
  ]
}
```

*Example output assumes the analysis date is November 9, 2026.*

## Project Structure

```text
invoice-api/
│
├── main.go
│
├── models/
│   └── invoice.go
│
├── handlers/
│   └── invoice.go
│
├── services/
│   └── invoice.go
│
├── routes/
│   └── routes.go
│
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

## File Responsibilities

| File | Responsibility |
|---|---|
| [`main.go`](main.go) | Application entry point. Initializes and starts the HTTP server. |
| [`models/invoice.go`](models/invoice.go) | Defines the invoice data structures and JSON fields. |
| [`handlers/invoice.go`](handlers/invoice.go) | Handles HTTP requests, parses input, and returns JSON responses. |
| [`services/invoice.go`](services/invoice.go) | Contains the core business logic: identifying overdue invoices and calculating overdue amounts and days. |
| [`routes/routes.go`](routes/routes.go) | Registers API endpoints and connects them to handlers. |
| [`.env`](.env) | Stores environment configuration. |
| [`.gitignore`](.gitignore) | Specifies files Git should ignore. |
| [`go.mod`](go.mod) | Defines the Go module and dependencies. |
| [`go.sum`](go.sum) | Verifies dependency integrity. |

## Request Flow

```text
       HTTP Request
            |
            v
         main.go
            |
            v
         routes/
            |
            v
        handlers/
            |
            v
        services/
            |
            v
      JSON Response
```

**Core responsibility:** Accept invoice data, identify overdue invoices, calculate the total overdue amount, and return the results.

**Note:** No database is required in the initial version. Each request is processed independently.
