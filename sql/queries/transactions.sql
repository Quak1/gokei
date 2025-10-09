-- name: CreateTransaction :one
INSERT INTO transactions (account_id, amount, category_id, title, attachment, note) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllTransactions :many
SELECT * FROM transactions;
