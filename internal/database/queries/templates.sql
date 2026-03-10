-- name: CreateTemplate :exec
INSERT INTO invoice_templates (id, name, path, is_default, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetTemplate :one
SELECT * FROM invoice_templates WHERE id = ?;

-- name: GetTemplateByPath :one
SELECT * FROM invoice_templates WHERE path = ?;

-- name: ListTemplates :many
SELECT * FROM invoice_templates ORDER BY name ASC;

-- name: GetDefaultTemplate :one
SELECT * FROM invoice_templates WHERE is_default = 1 LIMIT 1;

-- name: UpdateTemplate :exec
UPDATE invoice_templates SET name = ?, path = ?, is_default = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteTemplate :exec
DELETE FROM invoice_templates WHERE id = ?;
