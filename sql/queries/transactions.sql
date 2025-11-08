-- name: CreateTransaction :one
INSERT INTO transactions (account_id, amount_cents, category_id, title, attachment, note)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllTransactions :many
SELECT sqlc.embed(transactions) FROM transactions
INNER JOIN accounts ON transactions.account_id = accounts.id
WHERE accounts.user_id = $1;

-- name: GetTransactionsByAccountID :many
SELECT sqlc.embed(transactions) FROM transactions
INNER JOIN accounts ON transactions.account_id = accounts.id
WHERE transactions.account_id = $1 AND accounts.user_id = $2;

-- name: GetTransactionByID :one
SELECT sqlc.embed(transactions) FROM transactions
INNER JOIN accounts ON transactions.account_id = accounts.id
WHERE transactions.id = $1 AND accounts.user_id = $2;

-- name: DeleteTransactionByID :execresult
DELETE FROM transactions
USING accounts
WHERE transactions.account_id = accounts.id
  AND transactions.id = $1
  AND accounts.user_id = $2;

-- name: UpdateTransactionById :execresult
UPDATE transactions
SET amount_cents = $1,
    account_id = $2,
    category_id = $3,
    title = $4,
    date = $5,
    attachment = $6,
    note = $7,
    version = transactions.version + 1,
    updated_at = NOW()
FROM accounts
WHERE transactions.account_id = accounts.id
  AND transactions.id = $8
  AND accounts.user_id = $9
  AND transactions.version = $10;
