-- name: CreateSolution :one
INSERT INTO solutions (
    problemId, title, description
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetSolutions :many
SELECT id, problemId, title, description FROM solutions;
