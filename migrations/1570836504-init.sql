-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
