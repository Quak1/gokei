-- +goose Up
ALTER TABLE accounts
ADD user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE accounts
DROP COLUMN user_id;
