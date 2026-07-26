-- +goose Up
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_posts_created_at;