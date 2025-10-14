-- name: CreateAccount :one
INSERT INTO accounts (type, name) 
VALUES ($1, $2)
RETURNING *;

-- name: GetAllAccounts :many
SELECT * FROM accounts;

-- name: UpdateBalance :one
UPDATE accounts
SET balance_cents = balance_cents + $2
WHERE id = $1
RETURNING balance_cents;
