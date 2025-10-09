-- name: CreateCategory :one
INSERT INTO categories (name, color, icon) 
VALUES ($1, $2, $3)
RETURNING id, name, color, icon;

-- name: GetAllCategories :many
SELECT * FROM categories;
