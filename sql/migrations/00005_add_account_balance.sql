-- +goose Up
ALTER TABLE accounts
ADD balance_cents BIGINT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE accounts
DROP COLUMN balance_cents;
