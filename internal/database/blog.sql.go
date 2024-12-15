// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: blog.sql

package database

import (
	"context"
	"database/sql"
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

const createBlogImg = `-- name: CreateBlogImg :one
INSERT INTO blog_img(img_name, img, blog_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING img_name, img, blog_id, created_at, updated_at
`

type CreateBlogImgParams struct {
	ImgName   string
	Img       []byte
	BlogID    int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateBlogImg(ctx context.Context, arg CreateBlogImgParams) (BlogImg, error) {
	row := q.db.QueryRowContext(ctx, createBlogImg,
		arg.ImgName,
		arg.Img,
		arg.BlogID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i BlogImg
	err := row.Scan(
		&i.ImgName,
		&i.Img,
		&i.BlogID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteBlogImg = `-- name: DeleteBlogImg :one
DELETE FROM blog_img
WHERE blog_id = $1
RETURNING img_name, img, blog_id, created_at, updated_at
`

func (q *Queries) DeleteBlogImg(ctx context.Context, blogID int32) (BlogImg, error) {
	row := q.db.QueryRowContext(ctx, deleteBlogImg, blogID)
	var i BlogImg
	err := row.Scan(
		&i.ImgName,
		&i.Img,
		&i.BlogID,
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
WHERE b.is_published = TRUE
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
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username, bi.img_name
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
LEFT JOIN blog_img bi
    ON b.blog_id = bi.blog_id
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
	ImgName     sql.NullString
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
		&i.ImgName,
	)
	return i, err
}

const getBlogImg = `-- name: GetBlogImg :one
SELECT img_name, img, blog_id
FROM blog_img
WHERE blog_id = $1
`

type GetBlogImgRow struct {
	ImgName string
	Img     []byte
	BlogID  int32
}

func (q *Queries) GetBlogImg(ctx context.Context, blogID int32) (GetBlogImgRow, error) {
	row := q.db.QueryRowContext(ctx, getBlogImg, blogID)
	var i GetBlogImgRow
	err := row.Scan(&i.ImgName, &i.Img, &i.BlogID)
	return i, err
}

const getBlogsByCreator = `-- name: GetBlogsByCreator :many
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
WHERE b.created_by = $1
`

type GetBlogsByCreatorRow struct {
	BlogID      int32
	Title       string
	Body        []byte
	Excerpt     string
	IsPublished bool
	CreatedAt   time.Time
	Username    string
}

func (q *Queries) GetBlogsByCreator(ctx context.Context, createdBy int32) ([]GetBlogsByCreatorRow, error) {
	rows, err := q.db.QueryContext(ctx, getBlogsByCreator, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBlogsByCreatorRow
	for rows.Next() {
		var i GetBlogsByCreatorRow
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
