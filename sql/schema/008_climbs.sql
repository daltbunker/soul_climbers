
-- +goose Up
ALTER TABLE climb ADD COLUMN sub_areas TEXT;

-- +goose Down
ALTER TABLE climb DROP COLUMN sub_areas;