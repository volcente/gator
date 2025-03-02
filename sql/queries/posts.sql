-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description,
                   published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPostsForUser :many
SELECT p.title, p.published_at, p.description, p.url, p.feed_id
FROM posts AS p
INNER JOIN feeds as f ON p.feed_id = f.id
ORDER BY
  CASE WHEN @sort_by::varchar = 'title' AND @sort_order::varchar = 'asc' THEN p.title END ASC,
  CASE WHEN @sort_by= 'title' AND @sort_order = 'desc' THEN p.title END DESC,
  CASE WHEN @sort_by = 'url' AND @sort_order = 'asc' THEN p.url END ASC,
  CASE WHEN @sort_by = 'url' AND @sort_order = 'desc' THEN p.url END DESC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'asc' THEN p.published_at END ASC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'desc' THEN p.published_at END DESC,
  p.published_at DESC
LIMIT $1
OFFSET $2;