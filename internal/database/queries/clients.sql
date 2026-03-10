-- name: CreateClient :exec
INSERT INTO clients (id, name, email, phone, company, address, notes, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetClient :one
SELECT * FROM clients WHERE id = ?;

-- name: ListClients :many
SELECT * FROM clients ORDER BY name ASC;

-- name: UpdateClient :exec
UPDATE clients SET
    name = ?,
    email = ?,
    phone = ?,
    company = ?,
    address = ?,
    notes = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteClient :exec
DELETE FROM clients WHERE id = ?;

-- name: CountClients :one
SELECT COUNT(*) FROM clients;
