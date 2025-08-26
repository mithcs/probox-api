-- name: CreateUser :one
INSERT INTO users (
    username, password, full_name
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: GetUsers :many
SELECT id, username, full_name FROM users;

-- name: GetUserByUsername :one
SELECT id, username, full_name FROM users
WHERE username = $1 LIMIT 1;

-- name: GetCredentialsByUsername :one
SELECT id, username, password FROM users
WHERE username = $1 LIMIT 1;
