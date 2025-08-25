-- name: CreateUser :one
INSERT INTO users (
    username, password
) VALUES (
    $1, $2
)
RETURNING id;

-- name: GetUsers :many
SELECT id, username FROM users;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;
