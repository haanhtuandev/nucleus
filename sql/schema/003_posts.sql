-- +goose Up
ALTER TABLE posts
    ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN title TEXT NOT NULL,
    ADD COLUMN slug TEXT NOT NULL UNIQUE;

-- +goose Down
ALTER TABLE posts
    DROP COLUMN deleted_at,
    DROP COLUMN title,
    DROP COLUMN slug;