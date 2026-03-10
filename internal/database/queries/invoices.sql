-- name: CreateInvoice :exec
INSERT INTO invoices (id, invoice_number, client_id, status, issue_date, due_date, subtotal, tax_total, total, amount_paid, currency, notes, template_id, pdf_path, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = ?;

-- name: GetInvoiceByNumber :one
SELECT * FROM invoices WHERE invoice_number = ?;

-- name: ListInvoices :many
SELECT * FROM invoices ORDER BY created_at DESC;

-- name: ListInvoicesByStatus :many
SELECT * FROM invoices WHERE status = ? ORDER BY created_at DESC;

-- name: ListInvoicesByClient :many
SELECT * FROM invoices WHERE client_id = ? ORDER BY created_at DESC;

-- name: UpdateInvoiceStatus :exec
UPDATE invoices SET status = ?, updated_at = ? WHERE id = ?;

-- name: UpdateInvoiceTotals :exec
UPDATE invoices SET subtotal = ?, tax_total = ?, total = ?, updated_at = ? WHERE id = ?;

-- name: UpdateInvoiceAmountPaid :exec
UPDATE invoices SET amount_paid = ?, updated_at = ? WHERE id = ?;

-- name: UpdateInvoicePDFPath :exec
UPDATE invoices SET pdf_path = ?, updated_at = ? WHERE id = ?;

-- name: DeleteInvoice :exec
DELETE FROM invoices WHERE id = ?;
