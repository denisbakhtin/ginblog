-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE poststags(
  post_id INT REFERENCES posts(id) ON DELETE CASCADE,
  tag_name TEXT REFERENCES tags(name) ON DELETE CASCADE
);
CREATE UNIQUE INDEX posts_tags_idx ON poststags(post_id, tag_name);

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE poststags;

