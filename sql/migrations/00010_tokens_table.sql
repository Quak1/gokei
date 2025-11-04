-- +goose Up
CREATE TABLE tokens (
  hash bytea PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expiry TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE tokens;
