// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: blog.sql

package database

import (
	"context"
	"time"
)

const createBlog = `-- name: CreateBlog :one
INSERT INTO blog(body, title, excerpt, is_published, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING blog_id, body, title, excerpt, is_published, created_by, created_at, updated_at
`

type CreateBlogParams struct {
	Body        []byte
	Title       string
	Excerpt     string
	IsPublished bool
	CreatedBy   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreateBlog(ctx context.Context, arg CreateBlogParams) (Blog, error) {
	row := q.db.QueryRowContext(ctx, createBlog,
		arg.Body,
		arg.Title,
		arg.Excerpt,
		arg.IsPublished,
		arg.CreatedBy,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Blog
	err := row.Scan(
		&i.BlogID,
		&i.Body,
		&i.Title,
		&i.Excerpt,
		&i.IsPublished,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllBlogs = `-- name: GetAllBlogs :many
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
`

type GetAllBlogsRow struct {
	BlogID      int32
	Title       string
	Body        []byte
	Excerpt     string
	IsPublished bool
	CreatedAt   time.Time
	Username    string
}

func (q *Queries) GetAllBlogs(ctx context.Context) ([]GetAllBlogsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllBlogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllBlogsRow
	for rows.Next() {
		var i GetAllBlogsRow
		if err := rows.Scan(
			&i.BlogID,
			&i.Title,
			&i.Body,
			&i.Excerpt,
			&i.IsPublished,
			&i.CreatedAt,
			&i.Username,
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

const getBlogById = `-- name: GetBlogById :one
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username 
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
WHERE b.blog_id = $1
`

type GetBlogByIdRow struct {
	BlogID      int32
	Title       string
	Body        []byte
	Excerpt     string
	IsPublished bool
	CreatedAt   time.Time
	Username    string
}

func (q *Queries) GetBlogById(ctx context.Context, blogID int32) (GetBlogByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getBlogById, blogID)
	var i GetBlogByIdRow
	err := row.Scan(
		&i.BlogID,
		&i.Title,
		&i.Body,
		&i.Excerpt,
		&i.IsPublished,
		&i.CreatedAt,
		&i.Username,
	)
	return i, err
}

const updateBlog = `-- name: UpdateBlog :one
UPDATE blog
   SET title = $1,
       body = $2,
       excerpt = $3,
       is_published = $4,
       updated_at = $5
 WHERE blog_id = $6
RETURNING blog_id, body, title, excerpt, is_published, created_by, created_at, updated_at
`

type UpdateBlogParams struct {
	Title       string
	Body        []byte
	Excerpt     string
	IsPublished bool
	UpdatedAt   time.Time
	BlogID      int32
}

func (q *Queries) UpdateBlog(ctx context.Context, arg UpdateBlogParams) (Blog, error) {
	row := q.db.QueryRowContext(ctx, updateBlog,
		arg.Title,
		arg.Body,
		arg.Excerpt,
		arg.IsPublished,
		arg.UpdatedAt,
		arg.BlogID,
	)
	var i Blog
	err := row.Scan(
		&i.BlogID,
		&i.Body,
		&i.Title,
		&i.Excerpt,
		&i.IsPublished,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
