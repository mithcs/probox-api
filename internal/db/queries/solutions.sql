-- name: CreateSolution :one
INSERT INTO solutions (
    problem_id, title, description, owner_id, owner_name
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSolutions :many
SELECT * FROM solutions;

-- name: GetSolutionsByProblemID :many
SELECT * FROM solutions
WHERE problem_id = $1;

-- name: GetSolutionByID :one
SELECT * FROM solutions
WHERE id = $1 LIMIT 1;

-- name: GetSolutionCountByProblemID :one
SELECT count(*) FROM solutions
WHERE problem_id = $1;
