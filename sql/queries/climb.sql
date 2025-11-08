-- name: GetAllAreas :many
SELECT a.area_id, a.name, a.country, a.sub_areas
FROM area a;

-- name: GetArea :one
SELECT a.area_id, a.name, a.country, a.sub_areas
FROM area a
WHERE a.area_id = $1;

-- name: UpdateSubAreas :one
UPDATE area
SET sub_areas = $1
WHERE area_id = $2
RETURNING *;

-- name: CreateClimbDraft :one
INSERT INTO climb_draft(created_by, name, country, area_id, area, sub_areas, type)
VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT ON CONSTRAINT climb_draft_pkey
DO UPDATE SET 
    created_by = $1, 
    name = $2, 
    country = $3, 
    area_id = $4, 
    area = $5, 
    sub_areas = $6, 
    type = $7,
    updated_at = now() 
RETURNING *;

-- name: DeleteClimbDraft :one
DELETE FROM climb_draft
WHERE created_by = $1
RETURNING *;

-- name: GetClimbDraft :one
SELECT cd.created_by, cd.name, cd.country, cd.area_id, cd.area, cd.sub_areas, cd.type
FROM climb_draft cd
WHERE cd.created_by = $1;

-- name: CreateArea :one
INSERT INTO area(name, country, sub_areas, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateClimb :one
INSERT INTO climb(area_id, name, type, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllClimbs :many
SELECT c.area_id, c.name, c.type
FROM climb c
WHERE c.type = $1;

-- name: SearchClimbs :many
SELECT c.area_id, c.name, c.type
FROM climb c
WHERE c.type = $1
AND c.name % $2
ORDER BY similarity(c.name, $2) DESC, c.name;

-- name: SearchAreas :many
SELECT a.area_id, a.name, a.country, a.sub_areas
FROM area a
WHERE a.name % $1
ORDER BY similarity(a.name, $1) DESC, a.name;

-- name: SearchSubAreas :many
WITH sub_areas AS (
   SELECT a.area_id, a.name, unnest(string_to_array(a.sub_areas, ',')) as sub_area
   FROM area a ) 
SELECT sa.area_id, sa.name, sa.sub_area, similarity(sa.sub_area, $1) as sub_area_sml 
FROM sub_areas sa --WHERE sa.sub_area % $1
ORDER BY similarity(sa.sub_area, $1) DESC, sa.sub_area;
