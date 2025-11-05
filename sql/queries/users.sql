-- name: CreateUser :one
INSERT INTO users (username, name, password_hash) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserFromToken :one
SELECT users.id, users.username 
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.expiry > $2;

-- name: DeleteUserById :execresult
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserById :execresult
UPDATE users
SET name = $1, password_hash = $2, version = version + 1
WHERE id = $3 AND version = $4;
