-- +goose Up
ALTER TABLE accounts
ADD version INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE accounts
DROP COLUMN version;
