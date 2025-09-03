-- +goose Up
CREATE TYPE account_type AS ENUM ('credit', 'debit', 'cash');
CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind account_type NOT NULL,
  name TEXT NOT NULL
);

-- +goose Down
DROP TABLE accounts;
DROP TYPE account_type;
