-- +goose Up
CREATE TABLE users (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  username TEXT NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
