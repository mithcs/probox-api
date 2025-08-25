-- name: CreateProblem :one
INSERT INTO problems (
    title, description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetProblemById :one
SELECT * FROM problems
WHERE id = $1 LIMIT 1;
