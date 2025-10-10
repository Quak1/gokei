-- name: CreateCategory :one
INSERT INTO categories (name, color, icon) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllCategories :many
SELECT * FROM categories;
