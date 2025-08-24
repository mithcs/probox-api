-- name: CreateSolution :one
INSERT INTO solutions (
    problemId, solution
) VALUES (
    $1, $2
)
RETURNING *;
