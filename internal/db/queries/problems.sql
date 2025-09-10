-- name: CreateProblem :one
INSERT INTO problems (
    title, description, owner_id, owner_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetProblems :many
SELECT * FROM problems;

-- name: GetProblemByID :one
SELECT * FROM problems
WHERE id = $1 LIMIT 1;

-- name: DeleteProblemByID :exec
DELETE FROM problems WHERE id = $1;

-- name: SetAcceptedSolutionID :exec
UPDATE problems SET accepted_solution_id = $2
WHERE id = $1;

-- name: UnsetAcceptedSolutionID :exec
UPDATE problems SET accepted_solution_id = NULL
WHERE id = $1;
