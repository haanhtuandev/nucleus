-- +goose Up
ALTER TABLE posts
ADD COLUMN search_vector tsvector;

UPDATE posts
SET search_vector =
  to_tsvector('english', title || ' ' || content);

CREATE INDEX idx_posts_search
ON posts USING GIN(search_vector);

CREATE TRIGGER posts_search_vector_update
BEFORE INSERT OR UPDATE
ON posts
FOR EACH ROW
EXECUTE FUNCTION
tsvector_update_trigger(
  search_vector,
  'pg_catalog.english',
  title,
  content
);

-- +goose Down
-- 1. Remove the trigger first (it depends on the table and function)
DROP TRIGGER IF EXISTS posts_search_vector_update ON posts;

-- 2. Drop the index
DROP INDEX IF EXISTS idx_posts_search;

-- 3. Remove the column
-- This also removes the data stored in the tsvectors
ALTER TABLE posts
DROP COLUMN IF EXISTS search_vector;