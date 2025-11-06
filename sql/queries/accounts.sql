-- name: CreateAccount :one
INSERT INTO accounts (type, name, user_id) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllAccounts :many
SELECT * FROM accounts;

-- name: GetUserAccounts :many
SELECT * FROM accounts
WHERE user_id = $1;

-- name: UpdateBalance :one
UPDATE accounts
SET balance_cents = balance_cents + $1, updated_at = NOW()
WHERE id = $2 AND user_id = $3
RETURNING balance_cents;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1 AND user_id = $2;

-- name: DeleteAccountById :execresult
DELETE FROM accounts
WHERE id = $1 AND user_id = $2;

-- name: GetAccountSumBalance :one
SELECT accounts.name, SUM(transactions.amount_cents) AS balance
FROM transactions
RIGHT JOIN accounts ON transactions.account_id = accounts.id
WHERE account_id = $1 AND user_id = $2
GROUP BY accounts.id;

-- name: UpdateAccountById :execresult
UPDATE accounts
SET name = $1, type = $2, version = version + 1, updated_at = NOW()
WHERE id = $3 AND user_id = $4 AND version = $5;
