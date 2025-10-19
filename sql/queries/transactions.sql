-- name: CreateTransaction :one
INSERT INTO transactions (account_id, amount_cents, category_id, title, attachment, note)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllTransactions :many
SELECT * FROM transactions;

-- name: GetTransactionsByAccountID :many
SELECT * FROM transactions
WHERE account_id = $1;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = $1;

-- name: DeleteTransactionByID :execresult
DELETE FROM transactions
WHERE id = $1;

-- name: UpdateTransactionById :execresult
UPDATE transactions
SET amount_cents = $3, account_id = $4, category_id = $5, title = $6, date = $7, attachment = $8, note = $9, version = version + 1
WHERE id = $1 AND version = $2;
