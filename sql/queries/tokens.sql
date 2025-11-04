-- name: CreateToken :one
INSERT INTO tokens (hash, user_id, expiry)
VALUES ($1, $2, $3)
RETURNING *;
