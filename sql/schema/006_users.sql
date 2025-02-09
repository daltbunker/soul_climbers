-- +goose Up
ALTER TABLE users ADD COLUMN soul_score INT DEFAULT 0 NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN soul_score;