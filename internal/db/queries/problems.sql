-- name: CreateProblem :one
INSERT INTO problems (
    title, description, owner_id, owner_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetProblems :many
SELECT id, title, description, owner_id, owner_name, accepted_solution_id FROM problems;

-- name: GetProblemByID :one
SELECT id, title, description, owner_id, owner_name, accepted_solution_id FROM problems
WHERE id = $1 LIMIT 1;

-- name: DeleteProblemByID :exec
DELETE FROM problems WHERE id = $1;
