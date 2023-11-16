
-- +migrate Up
ALTER TABLE user ADD COLUMN created_at DATETIME NOT NULL;

-- +migrate Down
ALTER TABLE user REMOVE COLUMN created_at;
