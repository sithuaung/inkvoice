-- name: CreatePayment :exec
INSERT INTO payments (id, invoice_id, amount, method, reference, paid_at, notes, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListPaymentsByInvoice :many
SELECT * FROM payments WHERE invoice_id = ? ORDER BY paid_at DESC;

-- name: GetPayment :one
SELECT * FROM payments WHERE id = ?;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = ?;

-- name: SumPaymentsByInvoice :one
SELECT CAST(COALESCE(SUM(amount), 0) AS INTEGER) AS total_paid FROM payments WHERE invoice_id = ?;
