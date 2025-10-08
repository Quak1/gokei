-- +goose Up
CREATE TABLE transactions (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  amount INT NOT NULL,
  account_id INT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  date TIMESTAMP NOT NULL,
  attachment TEXT NOT NULL,
  note TEXT NOT NULL
);

-- +goose Down
DROP TABLE transactions;
