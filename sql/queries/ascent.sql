-- name: GetAscentsByClimb :many
SELECT a.grade, a.rating, TO_CHAR(a.ascent_date, 'YYYY-MM-DD') as ascent_date, a.over_200_pounds, a.attempts, a.comment, u.username
FROM ascent a
INNER JOIN users u
  ON u.users_id = a.created_by
WHERE a.climb_id = $1;

-- name: GetAscentsByUser :many
SELECT a.grade, a.rating, TO_CHAR(a.ascent_date, 'YYYY-MM-DD') as ascent_date, a.over_200_pounds, a.attempts, a.comment, u.username,
ar.area_id, ar.name, c.name as climb_name, c.climb_id, c.type
FROM ascent a
INNER JOIN users u
  ON u.users_id = a.created_by
INNER JOIN climb c
  ON c.climb_id = a.climb_id
INNER JOIN area ar
  ON ar.area_id = c.area_id
WHERE a.created_by = $1;

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

-- name: DeleteAscent :one
DELETE FROM ascent
WHERE created_by = $1
AND climb_id = $2
RETURNING *;
