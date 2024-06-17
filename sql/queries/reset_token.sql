-- name: CreateResetToken :one
INSERT INTO reset_token(token, expiration, email)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetResetTokenByToken :one
SELECT *
FROM reset_token 
WHERE token = $1;