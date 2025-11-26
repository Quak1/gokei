-- name: CreateRecurringTransaction :one
INSERT INTO recurring_transactions (
    account_id,
    amount_cents,
    category_id,
    title,
    note,
    frequency,
    interval,
    start_date,
    end_date,
    day_month,
    day_week,
    max_occurrences,
    is_active
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetUserRecurringTransactions :many
SELECT rt.* FROM recurring_transactions rt
INNER JOIN accounts ON rt.account_id = accounts.id
WHERE accounts.user_id = $1;

-- name: GetRecurringTransactionByID :one
SELECT rt.* FROM recurring_transactions rt
INNER JOIN accounts ON rt.account_id = accounts.id
WHERE rt.id = $1 AND accounts.user_id = $2;

-- name: GetActiveRecurringTransactions :many
SELECT rt.* FROM recurring_transactions rt
INNER JOIN accounts ON rt.account_id = accounts.id
WHERE accounts.user_id = $1 
  AND rt.is_active = true
  AND rt.start_date <= $2
  AND (rt.end_date IS NULL OR rt.end_date >= $2);

-- name: UpdateRecurringTransaction :execresult
UPDATE recurring_transactions
SET amount_cents = $1,
    category_id = $2,
    title = $3,
    note = $4,
    frequency = $5,
    interval = $6,
    end_date = $7,
    day_month = $8,
    day_week = $9,
    max_occurrences = $10,
    is_active = $11,
    version = version + 1,
    updated_at = NOW()
FROM accounts
WHERE recurring_transactions.account_id = accounts.id
  AND recurring_transactions.id = $12
  AND accounts.user_id = $13
  AND recurring_transactions.version = $14;

-- name: DeleteRecurringTransaction :execresult
DELETE FROM recurring_transactions
USING accounts
WHERE recurring_transactions.account_id = accounts.id
  AND recurring_transactions.id = $1
  AND accounts.user_id = $2;

-- name: CreateOccurrence :one
INSERT INTO recurring_transaction_occurrences (
    recurring_transaction_id,
    transaction_id,
    occurrence_date
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOccurrences :many
SELECT * FROM recurring_transaction_occurrences
WHERE recurring_transaction_id = $1;

-- name: GetOccurrenceForDate :one
SELECT * FROM recurring_transaction_occurrences
WHERE recurring_transaction_id = $1 AND occurrence_date = $2;

-- name: GetLastOccurrence :one
SELECT * FROM recurring_transaction_occurrences
WHERE recurring_transaction_id = $1
ORDER BY occurrence_date DESC
LIMIT 1;
