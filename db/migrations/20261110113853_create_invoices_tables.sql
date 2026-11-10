-- +goose Up

CREATE TABLE invoices (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    customer VARCHAR(50) NOT NULL,

    invoice_no VARCHAR(50) NOT NULL UNIQUE,

    required_amount BIGINT NOT NULL
        CHECK (required_amount >= 0),

    paid_amount BIGINT NOT NULL DEFAULT 0
        CHECK (paid_amount >= 0),

    due_date DATE NOT NULL,

    transaction_status VARCHAR(20) NOT NULL,

    transaction_completion_date TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_invoice_payment
        CHECK (paid_amount <= required_amount)
);


-- +goose Down

DROP TABLE invoices;


-- TODO: NOTE: FIX: Up describes the change you want to apply.
--
-- Down describes how to reverse that change.
