-- name: CreateSolution :one
INSERT INTO solutions (
    problemId, title, description
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetSolutions :many
SELECT id, problemId, title, description FROM solutions;

-- name: GetSolutionsByProblemId :many
SELECT id, problemId, title, description FROM solutions
WHERE problemId = $1;
