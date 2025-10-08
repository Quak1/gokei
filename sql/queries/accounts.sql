-- name: CreateAccount :one
INSERT INTO accounts (type, name) 
VALUES ($1, $2)
RETURNING id, type, name;
