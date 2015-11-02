-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tags(
  name TEXT PRIMARY KEY
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE tags;

