
-- +migrate Up
CREATE TABLE user (
    user_id     CHAR(36)        PRIMARY KEY,
    username    VARCHAR(256)    NOT NULL
);

-- +migrate Down
DROP TABLE user;
