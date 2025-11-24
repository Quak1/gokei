-- name: CreateCategory :one
INSERT INTO categories (user_id, name, color, icon) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllCategories :many
SELECT * FROM categories
WHERE user_id = @admin_id OR user_id = @user_id;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 AND user_id = $2;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1 AND user_id = $2;

-- name: DeleteCategoryById :execresult
DELETE FROM categories
WHERE id = $1 AND user_id = $2;

-- name: UpdateCategoryById :execresult
UPDATE categories
SET name = $1, color = $2, icon = $3, version = version + 1, updated_at = NOW()
WHERE id = $4 AND user_id = $5 AND version = $6;
