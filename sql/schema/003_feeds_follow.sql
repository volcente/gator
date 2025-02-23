-- +goose Up
CREATE TABLE feed_follows (
  id uuid PRIMARY KEY,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  user_id uuid NOT NULL,
  feed_id uuid NOT NULL,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
  UNIQUE (user_id, feed_id)
);

ALTER TABLE feeds
RENAME COLUMN user_id to author_id;

-- +goose Down
DROP TABLE feed_follows;

ALTER TABLE feeds
RENAME COLUMN author_id to user_id;
