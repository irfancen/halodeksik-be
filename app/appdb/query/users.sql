-- name: CreateUser :one
INSERT INTO users(email, password, user_role_id, is_verified)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;