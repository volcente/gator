-- name: CreateFeed :one
INSERT INTO
  feeds (id, created_at, updated_at, name, url, author_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1 LIMIT 1;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name as username
FROM feeds as f
INNER JOIN users as u
ON f.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT u.name as username, f.name as feed_name
FROM users as u
INNER JOIN feed_follows as ff ON ff.user_id = u.id
INNER JOIN feeds as f ON f.id = ff.feed_id
WHERE u.id = $1;

-- name: CreateFeedFollow :one
WITH inserted_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT ff.*, f.name, u.name as username
FROM inserted_follow as ff
INNER JOIN feeds as f ON f.id = ff.feed_id
INNER JOIN users as u ON u.id = ff.user_id;

-- name: DeleteFeedFollower :exec
WITH feed AS (
  SELECT f.url FROM feeds as f
  WHERE f.url = $1
)
DELETE FROM feed_follows as ff
USING feed as f
WHERE ff.user_id = $2;
