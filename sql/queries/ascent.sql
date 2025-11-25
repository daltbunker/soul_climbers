-- name: GetAscentsByClimb :many
SELECT a.grade, a.rating, TO_CHAR(a.ascent_date, 'YYYY-MM-DD') as ascent_date, a.over_200_pounds, a.attempts, a.comment, u.username
FROM ascent a
INNER JOIN users u
  ON u.users_id = a.created_by
WHERE a.climb_id = $1;

-- name: CreateOrUpdateAscent :one
INSERT INTO ascent(grade, rating, ascent_date, over_200_pounds, attempts, comment, created_by, climb_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
ON CONFLICT ON CONSTRAINT ascent_pkey 
DO UPDATE SET
    grade = $1,
    rating = $2,
    ascent_date = $3,
    over_200_pounds = $4,
    attempts = $5,
    comment = $6,
    created_by = $7,
    updated_at = now() 
RETURNING *;