-- name: CreateBlog :one
INSERT INTO blog(body, title, excerpt, is_published, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBlogById :one
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username 
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
WHERE b.blog_id = $1;

-- name: GetAllBlogs :many
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id;

-- name: UpdateBlog :one
UPDATE blog
   SET title = $1,
       body = $2,
       excerpt = $3,
       is_published = $4,
       updated_at = $5
 WHERE blog_id = $6
RETURNING *;