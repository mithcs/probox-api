-- name: CreateProblem :one
INSERT INTO problems (
    title, description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetProblems :many
SELECT id, title, description FROM problems;

-- name: GetProblemById :one
SELECT * FROM problems
WHERE id = $1 LIMIT 1;

