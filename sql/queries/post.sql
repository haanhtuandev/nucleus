-- name: CreatePost :one
INSERT INTO posts (id, title, content, user_id, created_at, updated_at, deleted_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW(),
    NULL
)
RETURNING *;

-- name: GetAllPosts :many
SELECT * FROM posts WHERE deleted_at is NULL;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = $1 AND deleted_at is NULL;

-- name: DeletePost :exec
UPDATE posts
SET deleted_at = NOW()
WHERE id = $1;



-- name: GetPostsByUser :many
SELECT * from posts
WHERE user_id = $1 AND deleted_at is NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePostInfo :exec
UPDATE posts
SET title = $1, content = $2, updated_at = NOW()
WHERE id = $3;


-- name: FetchPost :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPostsCount :one
SELECT COUNT(*) FROM posts
WHERE user_id = $1;

-- name: SearchPosts :many
SELECT
  posts.title,
  posts.content,
  posts.created_at,
  posts.updated_at,
  users.username,
  users.id
FROM posts
JOIN users ON posts.user_id = users.id
WHERE posts.search_vector @@ plainto_tsquery('english', $1)
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchPosts :one
SELECT COUNT(*)
FROM posts
WHERE search_vector @@ plainto_tsquery('english', $1);