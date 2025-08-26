-- name: CreateSolution :one
INSERT INTO solutions (
    problem_id, title, description
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetSolutions :many
SELECT id, problem_id, title, description FROM solutions;

-- name: GetSolutionsByProblemID :many
SELECT id, problem_id, title, description FROM solutions
WHERE problem_id = $1;

-- name: GetSolutionByID :one
SELECT id, problem_id, title, description FROM solutions
WHERE id = $1 LIMIT 1;
