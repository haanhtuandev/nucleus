-- name: CreatePost :one
INSERT INTO posts (id, title, content, user_id, slug, created_at, updated_at, deleted_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
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


-- name: GetPostBySlug :one
SELECT * FROM posts WHERE slug = $1 AND deleted_at is NULL;


-- name: LookUpSlug :one
SELECT COUNT(*) from posts WHERE slug = $1;

-- name: GetPostsByUser :many
SELECT posts.id, posts.title, posts.created_at, posts.updated_at 
FROM posts JOIN users ON posts.user_id = users.id
WHERE deleted_at is NULL AND posts.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePostInfo :exec
UPDATE posts
SET title = $1, content = $2, slug = $3, updated_at = NOW()
WHERE id = $4;


-- name: FetchPost :many
SELECT posts.title, posts.content, posts.created_at, posts.updated_at, users.username, users.id 
FROM posts JOIN users on posts.user_id = users.id
ORDER BY posts.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPostsCount :one
SELECT COUNT(*) FROM posts JOIN users on posts.user_id = users.id
WHERE users.id = $1;

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