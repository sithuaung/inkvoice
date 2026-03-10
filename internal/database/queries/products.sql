-- name: CreateProduct :exec
INSERT INTO products (id, name, description, unit_price, currency, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetProduct :one
SELECT * FROM products WHERE id = ?;

-- name: ListProducts :many
SELECT * FROM products ORDER BY name ASC;

-- name: UpdateProduct :exec
UPDATE products SET
    name = ?,
    description = ?,
    unit_price = ?,
    currency = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;
