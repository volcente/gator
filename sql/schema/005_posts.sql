-- +goose Up
CREATE TABLE posts (
  id           uuid PRIMARY KEY,
  created_at   timestamp NOT NULL,
  updated_at   timestamp NOT NULL,
  title        text      NOT NULL,
  url          text      NOT NULL UNIQUE,
  description  text      NOT NULL,
  published_at timestamp NOT NULL,
  feed_id      uuid      NOT NULL,
  CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- +goose Down
DROP TABLE posts;