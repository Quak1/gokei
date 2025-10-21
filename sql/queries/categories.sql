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
SET name = $1, color = $2, icon = $3, version = version + 1
WHERE id = $4 AND version = $5;

-- name: InsertInitialCategory :exec
INSERT INTO categories (name, color, icon)
SELECT 'InitialBalance', '#FFF', 'icon'
WHERE NOT EXISTS (SELECT 1 FROM categories);
