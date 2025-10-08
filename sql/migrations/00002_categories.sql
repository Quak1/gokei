-- +goose Up
CREATE TABLE categories (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  name TEXT NOT NULL,
  color TEXT NOT NULL,
  icon TEXT NOT NULl
);

-- +goose Down
DROP TABLE categories;
