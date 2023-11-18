
-- +migrate Up
ALTER TABLE shopping_list ADD COLUMN item_count INTEGER NOT NULL;

-- +migrate Down
ALTER TABLE shopping_list REMOVE COLUMN item_count;
