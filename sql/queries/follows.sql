-- name: CreateFollow :exec
INSERT INTO follows (followee_id, follower_id, created_at)
VALUES (
    $1,
    $2,
    NOW()
);

-- name: DeleteFollow :exec
DELETE FROM follows
WHERE followee_id = $1 AND follower_id = $2;

-- name: GetFollowers :many
SELECT COUNT(*), follower_id FROM users JOIN follows on users.id = follows.followee_id
WHERE users.id = $1
ORDER BY follows.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetFollowees :many
SELECT COUNT(*), followee_id FROM users JOIN follows on users.id = follows.follower_id
WHERE users.id = $1
ORDER BY follows.created_at DESC
LIMIT $2 OFFSET $3;