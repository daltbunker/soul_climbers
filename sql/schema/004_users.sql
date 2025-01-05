-- +goose Up
ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'user' NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN role;