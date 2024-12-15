-- name: CreateBlog :one
INSERT INTO blog(body, title, excerpt, is_published, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBlogById :one
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username, bi.img_name
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
LEFT JOIN blog_img bi
    ON b.blog_id = bi.blog_id
WHERE b.blog_id = $1;

-- name: GetAllBlogs :many
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
WHERE b.is_published = TRUE;

-- name: UpdateBlog :one
UPDATE blog
   SET title = $1,
       body = $2,
       excerpt = $3,
       is_published = $4,
       updated_at = $5
 WHERE blog_id = $6
RETURNING *;

-- name: CreateBlogImg :one 
INSERT INTO blog_img(img_name, img, blog_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBlogImg :one
SELECT img_name, img, blog_id
FROM blog_img
WHERE blog_id = $1;

-- name: DeleteBlogImg :one
DELETE FROM blog_img
WHERE blog_id = $1
RETURNING *;

-- name: GetBlogsByCreator :many
SELECT b.blog_id, b.title, b.body, b.excerpt, b.is_published, b.created_at, u.username
FROM blog b
INNER JOIN users u
    ON b.created_by = u.users_id
WHERE b.created_by = $1;