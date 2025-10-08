-- +goose Up
CREATE TYPE account_type AS ENUM ('credit', 'debit', 'cash');
CREATE TABLE accounts (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  type account_type NOT NULL,
  name TEXT NOT NULL
);

-- +goose Down
DROP TABLE accounts;
DROP TYPE account_type;
