
-- +migrate Up
CREATE TABLE shopping_list (
    list_id         CHAR(36)        PRIMARY KEY,
    creator_id      CHAR(36)        NOT NULL,
    list_name       VARCHAR(1000)   NOT NULL,

    FOREIGN KEY (creator_id) REFERENCES user (user_id)
);

CREATE TABLE shopping_list_item (
    item_id     CHAR(36)        PRIMARY KEY,
    item_name   VARCHAR(1000)   NOT NULL,
    list_id     CHAR(36)        NOT NULL,

    FOREIGN KEY (list_id)   REFERENCES shopping_list (list_id)
);

-- +migrate Down
DROP TABLE shopping_list_item;
DROP TABLE shopping_list;
