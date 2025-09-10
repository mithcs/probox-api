-- name: CreateUser :one
INSERT INTO users (
    username, password, full_name
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: GetUsers :many
SELECT id, username, full_name, created_at FROM users;

-- name: GetUserByID :one
SELECT id, username, full_name, created_at FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, full_name, created_at FROM users
WHERE username = $1 LIMIT 1;

-- name: GetCredentialsByUsername :one
SELECT id, username, password FROM users
WHERE username = $1 LIMIT 1;

-- name: GetFullNameByID :one
SELECT full_name FROM users
WHERE id = $1 LIMIT 1;

-- name: DeleteUserByID :exec
UPDATE users SET
    username = '_user' || id,
    full_name = 'Deleted User',
    password = $2
WHERE id = $1;
