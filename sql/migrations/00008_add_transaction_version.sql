-- +goose Up
ALTER TABLE transactions
ADD version INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE transactions
DROP COLUMN version;
