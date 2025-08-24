-- name: CreateProblem :one
INSERT INTO problems (
    title, description
) VALUES (
    $1, $2
)
RETURNING id;
