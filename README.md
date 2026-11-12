# B2B Invoice API

A REST API built with **Go, Gin and PostgreSQL** for creating and storing B2B invoices.

Supports multiple invoice submissions in a single request, with payment validation and automatic transaction status calculation.

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL
- **Driver:** pgx
- **Deployment:** Render

## Features

- Create single or multiple invoices.
- Validate invoice and payment data.
- Prevent overpayments.
- Automatically calculate payment status (`pending` / `fulfilled`).
- Record payment completion timestamps.
- Use database transactions for atomic batch insertion.

## API

**Base URL:** https://b2b-invoice.onrender.com

### Create Invoices

`POST /invoices`

Request body:

```json
[
  {
    "customer": "ABC Traders",
    "invoice_no": "INV-001",
    "required_amount": 15000,
    "paid_amount": 5000,
    "due_date": "2026-11-30"
  },
  {
    "customer": "XYZ Stores",
    "invoice_no": "INV-002",
    "required_amount": 8000,
    "paid_amount": 8000,
    "due_date": "2026-12-10"
  }
]
```

Response (`201 Created`):

```json
{
  "message": "Invoices created successfully",
  "count": 2
}
```

### Validation

- Customer and invoice number are required.
- Invoice numbers must be unique.
- Required amount must be greater than zero.
- Paid amount cannot exceed the required amount.
- Due dates must follow `YYYY-MM-DD`.
- Empty invoice batches are rejected.

## Project Structure

```text
.
├── db/
│   ├── database.go
│   └── migrations/
├── handlers/
│   └── invoiceHandler.go
├── models/
│   └── invoiceModel.go
├── routes/
│   └── invoiceRoutes.go
├── services/
├── main.go
├── go.mod
└── go.sum
```

## Run Locally

Clone the repository:

```bash
git clone https://github.com/Anjas2005/B2B_INVOICE.git
cd B2B_INVOICE
```

Install dependencies:

```bash
go mod download
```

Configure environment variables in `.env`:

```env
DATABASE_URL=your_postgresql_connection_string
GIN_MODE=release
PORT=8080
```

Apply the SQL migration in `db/migrations/` to your PostgreSQL database. Execute the `Up` section only.

Start the server:

```bash
go run .
```

The API will run on `http://localhost:8080`.

## Deployment

Deployed on Render:

https://b2b-invoice.onrender.com

**Build command:**

```bash
go build -tags netgo -ldflags '-s -w' -o app .
```

**Start command:**

```bash
./app
```

Configure `DATABASE_URL` in Render's environment settings.

---

**Note:** All invoices in a batch are inserted within a single database transaction. If an insertion fails, the transaction is rolled back.
