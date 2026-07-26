-- +goose Up
ALTER TABLE posts
DROP COLUMN slug;


-- +goose Down
ALTER TABLE posts
ADD COLUMN slug TEXT NOT NULL UNIQUE;