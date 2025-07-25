-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, 
    $2
)
RETURNING *;

-- name: FindUser :one
SELECT id
FROM users
WHERE email = $1;

-- name: CheckPassword :one
SELECT id
FROM users
WHERE password = $1;