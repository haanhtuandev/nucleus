-- name: CreateUser :one
INSERT INTO users (id, username, hashed_password, bio, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
RETURNING *;


-- name: GetAllUsers :many
SELECT * FROM users;


-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByName :one
SELECT * FROM users WHERE username = $1;

-- name: RefreshDB :exec
TRUNCATE TABLE users cascade;