-- name: CreateRecurringInvoice :exec
INSERT INTO recurring_invoices (id, client_id, schedule, status, next_run, last_run, currency, template_id, notes, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetRecurringInvoice :one
SELECT * FROM recurring_invoices WHERE id = ?;

-- name: ListRecurringInvoices :many
SELECT * FROM recurring_invoices ORDER BY created_at DESC;

-- name: UpdateRecurringInvoiceStatus :exec
UPDATE recurring_invoices SET status = ?, updated_at = ? WHERE id = ?;

-- name: UpdateRecurringInvoiceNextRun :exec
UPDATE recurring_invoices SET next_run = ?, last_run = ?, updated_at = ? WHERE id = ?;

-- name: ListDueRecurringInvoices :many
SELECT * FROM recurring_invoices WHERE status = 'active' AND next_run <= ? ORDER BY next_run ASC;

-- name: DeleteRecurringInvoice :exec
DELETE FROM recurring_invoices WHERE id = ?;

-- name: CreateRecurringInvoiceItem :exec
INSERT INTO recurring_invoice_items (id, recurring_invoice_id, product_id, description, quantity, unit_price, tax_id, tax_rate, sort_order, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListRecurringInvoiceItems :many
SELECT * FROM recurring_invoice_items WHERE recurring_invoice_id = ? ORDER BY sort_order ASC;

-- name: DeleteRecurringInvoiceItems :exec
DELETE FROM recurring_invoice_items WHERE recurring_invoice_id = ?;
