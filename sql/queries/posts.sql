-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description,
                   published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPostsForUser :many
SELECT p.title, p.published_at, p.description, p.url, p.feed_id
FROM posts AS p
INNER JOIN feeds as f ON p.feed_id = f.id
ORDER BY p.published_at DESC
LIMIT $1
OFFSET $2;
