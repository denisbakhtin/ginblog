-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE pages(
  id SERIAL,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  published boolean NOT NULL DEFAULT true
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE pages;

