-- +goose Up
CREATE TYPE recurrence_frequency AS ENUM ('daily', 'weekly', 'monthly', 'yearly');

CREATE TABLE recurring_transactions (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  version INT NOT NULL DEFAULT 1,

  account_id INT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  amount_cents INT NOT NULL,
  category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  note TEXT NOT NULL,

  frequency recurrence_frequency NOT NULL,
  interval INT NOT NULL DEFAULT 1,
  start_date TIMESTAMP NOT NULL,
  end_date TIMESTAMP, 

  day_month INT,
  day_week INT,
  max_occurrences INT,

  is_active BOOLEAN NOT NULL DEFAULT true,

  CONSTRAINT valid_interval CHECK (interval > 0),
  CONSTRAINT valid_day_of_month CHECK (day_month IS NULL OR (day_month >= 1 AND day_month <= 31)),
  CONSTRAINT valid_day_of_week CHECK (day_week IS NULL OR (day_week >= 0 AND day_week <= 6))
);

CREATE TABLE recurring_transaction_occurrences (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    recurring_transaction_id INT NOT NULL REFERENCES recurring_transactions(id) ON DELETE CASCADE,
    transaction_id INT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    occurrence_date DATE NOT NULL,

    UNIQUE(recurring_transaction_id, occurrence_date)
);

CREATE INDEX idx_recurring_transactions_account ON recurring_transactions(account_id);
CREATE INDEX idx_recurring_transactions_active ON recurring_transactions(is_active) WHERE is_active = true;
CREATE INDEX idx_occurrences_recurring ON recurring_transaction_occurrences(recurring_transaction_id);

-- +goose Down
DROP TABLE recurring_transaction_occurrences;
DROP TABLE recurring_transactions;
DROP TYPE recurrence_frequency;
