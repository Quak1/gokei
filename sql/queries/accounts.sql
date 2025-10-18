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

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1;

-- name: DeleteAccountById :execresult
DELETE FROM accounts
WHERE id = $1;

-- name: GetAccountSumBalance :one
SELECT accounts.name, SUM(transactions.amount_cents) AS balance
FROM transactions
RIGHT JOIN accounts ON transactions.account_id = accounts.id
WHERE account_id = $1
GROUP BY accounts.id;
