-- name: CreateUser :exec
INSERT INTO users (id, email, username, hashedPassword, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CheckUser :one
SELECT id, email FROM users WHERE email = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;