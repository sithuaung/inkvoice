-- name: CreateInvoiceItem :exec
INSERT INTO invoice_items (id, invoice_id, product_id, description, quantity, unit_price, tax_id, tax_rate, amount, sort_order, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListInvoiceItems :many
SELECT * FROM invoice_items WHERE invoice_id = ? ORDER BY sort_order ASC;

-- name: DeleteInvoiceItem :exec
DELETE FROM invoice_items WHERE id = ?;

-- name: DeleteInvoiceItemsByInvoice :exec
DELETE FROM invoice_items WHERE invoice_id = ?;

-- name: GetInvoiceItem :one
SELECT * FROM invoice_items WHERE id = ?;
