-- +goose Up
ALTER TABLE transactions 
ALTER COLUMN amount TYPE BIGINT;

ALTER TABLE transactions
RENAME COLUMN amount TO amount_cents;

-- +goose Down
ALTER TABLE transactions
RENAME COLUMN amount_cents TO amount;

ALTER TABLE transactions;
ALTER COLUMN amount TYPE INT;
