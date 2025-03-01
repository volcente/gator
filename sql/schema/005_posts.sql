-- +goose Up
-- id - a unique identifier for the post
-- created_at - the time the record was created
-- updated_at - the time the record was last updated
-- title - the title of the post
-- url - the URL of the post (this should be unique)
-- description - the description of the post
-- published_at - the time the post was published
-- feed_id - the ID of the feed that the post came from
CREATE TABLE posts
(
    id           uuid PRIMARY KEY,
    created_at   timestamp NOT NULL,
    updated_at   timestamp NOT NULL,
    title        text      NOT NULL,
    url          text      NOT NULL UNIQUE,
    description  text,
    published_at timestamp NOT NULL,
    feed_id      uuid      NOT NULL,
    CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- +goose Down
DROP TABLE posts;