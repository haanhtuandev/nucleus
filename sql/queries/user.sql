-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW()
)
RETURNING *;