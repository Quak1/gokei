-- +goose Up
CREATE TABLE users (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  version INT NOT NULL DEFAULT 1,
  username TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password_hash bytea NOT NULL
);

-- +goose Down
DROP TABLE users;
