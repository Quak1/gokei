-- name: CreateCategory :one
INSERT INTO categories (name, color, icon) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: DeleteCategoryById :execresult
DELETE FROM categories
WHERE id = $1;

-- name: UpdateCategoryById :execresult
UPDATE categories
SET name = $1, color = $2, icon = $3, version = version + 1, updated_at = NOW()
WHERE id = $4 AND version = $5;
