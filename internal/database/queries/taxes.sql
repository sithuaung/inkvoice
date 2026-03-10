-- name: CreateTax :exec
INSERT INTO taxes (id, name, rate, is_default, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetTax :one
SELECT * FROM taxes WHERE id = ?;

-- name: ListTaxes :many
SELECT * FROM taxes ORDER BY name ASC;

-- name: UpdateTax :exec
UPDATE taxes SET name = ?, rate = ?, is_default = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteTax :exec
DELETE FROM taxes WHERE id = ?;

-- name: GetDefaultTax :one
SELECT * FROM taxes WHERE is_default = 1 LIMIT 1;
