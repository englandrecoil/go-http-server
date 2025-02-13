-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN DEFAULT false NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;