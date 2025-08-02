-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, password, nickname)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, 
    $2,
    NULL
)
RETURNING *;

-- name: FindUser :one
SELECT id
FROM users
WHERE username = $1;

-- name: SetNickname :one
INSERT INTO users (nickname)
VALUES ($1)
RETURNING *;

-- name: CheckPassword :one
SELECT id, username, created_at, updated_at, nickname
FROM users
WHERE password = $1;