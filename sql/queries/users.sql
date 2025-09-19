-- name: CreateUser :one
INSERT INTO users (username, hashed_password, name) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;
