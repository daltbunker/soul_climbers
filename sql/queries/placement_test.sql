-- name: CreateTestQuestion :one
INSERT INTO test_question(question_text, input_type, answers, answer_points, points, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTestQuestions :many
SELECT tq.test_question_id, tq.question_text, tq.input_type, tq.answers, tq.answer_points, tq.points, tq.created_at
FROM test_question tq;

-- name: InsertPlacmentTest :one
INSERT INTO placement_test(user_id, score, created_at, updated_at)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetPlacementTest :one
SELECT pt.score 
FROM placement_test pt
WHERE pt.user_id = $1;