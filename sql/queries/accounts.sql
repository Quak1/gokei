-- name: CreateAccount :one
INSERT INTO accounts (type, name) 
VALUES ($1, $2)
RETURNING *;

-- name: GetAllAccounts :many
SELECT * FROM accounts;
