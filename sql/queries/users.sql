-- name: CreateUser :one
INSERT INTO users(username, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: SetUserSoulScore :one
UPDATE users
SET soul_score = $1 
WHERE username = $2
RETURNING *;