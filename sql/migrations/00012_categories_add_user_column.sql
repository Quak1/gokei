-- +goose Up
ALTER TABLE categories
ADD user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE categories
DROP COLUMN user_id;
