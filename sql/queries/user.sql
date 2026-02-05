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


-- name: GetAllUsers :many
SELECT * FROM users;


-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: RefreshDB :exec
TRUNCATE TABLE users cascade;