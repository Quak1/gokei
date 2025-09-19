-- +goose Up
CREATE TABLE categories (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  color TEXT NOT NULL,
  icon TEXT NOT NULl
);

-- +goose Down
DROP TABLE categories;
