-- name: CreateBlog :one
INSERT INTO blog(body, title, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBlogById :one
SELECT b.blog_id, b.title, b.body, b.created_at, u.username 
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
WHERE b.blog_id = $1;

-- name: GetAllBlogs :many
SELECT b.blog_id, b.title, b.body, b.created_at, u.username
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id;