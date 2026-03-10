-- name: GetSettings :one
SELECT * FROM settings WHERE id = 'default';

-- name: UpdateSettings :exec
UPDATE settings SET
    company_name = ?,
    company_email = ?,
    company_phone = ?,
    company_address = ?,
    invoice_prefix = ?,
    next_invoice_number = ?,
    default_due_days = ?,
    default_currency = ?,
    default_template_id = ?,
    updated_at = ?
WHERE id = 'default';

-- name: IncrementInvoiceNumber :one
UPDATE settings SET next_invoice_number = next_invoice_number + 1, updated_at = ?
WHERE id = 'default'
RETURNING next_invoice_number - 1;
