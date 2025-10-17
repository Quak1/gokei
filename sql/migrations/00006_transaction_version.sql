-- +goose Up
ALTER TABLE categories
ADD version INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE categories
DROP COLUMN version;
