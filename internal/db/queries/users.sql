-- name: CreateUser :one
INSERT INTO users (
    username, password
) VALUES (
    $1, $2
)
RETURNING id;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;
