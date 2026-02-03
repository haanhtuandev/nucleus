-- name: CreateUser :one
INSERT INTO users (id, username, bio, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
RETURNING *;