-- name: CreateAccount :one
INSERT INTO accounts (owner_id, kind, name) 
VALUES ($1, $2, $3)
RETURNING *;
