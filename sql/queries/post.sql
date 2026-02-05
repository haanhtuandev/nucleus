-- name: CreatePost :one
INSERT INTO posts (id, content, user_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetAllPosts :many
SELECT * FROM posts;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = $1;