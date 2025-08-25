-- name: CreateSolution :one
INSERT INTO solutions (
    problemId, solution
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetSolutions :many
SELECT id, problemId, solution FROM solutions;
