// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description,
                   published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.ExecContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	return err
}

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT p.title, p.published_at, p.description, p.url, p.feed_id
FROM posts AS p
INNER JOIN feeds as f ON p.feed_id = f.id
ORDER BY
  CASE WHEN $3::varchar = 'title' AND $4::varchar = 'asc' THEN p.title END ASC,
  CASE WHEN $3= 'title' AND $4 = 'desc' THEN p.title END DESC,
  CASE WHEN $3 = 'url' AND $4 = 'asc' THEN p.url END ASC,
  CASE WHEN $3 = 'url' AND $4 = 'desc' THEN p.url END DESC,
  CASE WHEN $3 = 'published_at' AND $4 = 'asc' THEN p.published_at END ASC,
  CASE WHEN $3 = 'published_at' AND $4 = 'desc' THEN p.published_at END DESC,
  p.published_at DESC
LIMIT $1
OFFSET $2
`

type GetPostsForUserParams struct {
	Limit     int32
	Offset    int32
	SortBy    string
	SortOrder string
}

type GetPostsForUserRow struct {
	Title       string
	PublishedAt time.Time
	Description string
	Url         string
	FeedID      uuid.UUID
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]GetPostsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser,
		arg.Limit,
		arg.Offset,
		arg.SortBy,
		arg.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsForUserRow
	for rows.Next() {
		var i GetPostsForUserRow
		if err := rows.Scan(
			&i.Title,
			&i.PublishedAt,
			&i.Description,
			&i.Url,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
