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

-- name: FindID :one
SELECT id
FROM users
WHERE username = $1;

-- name: SetNickname :one
UPDATE users
SET nickname = $1
WHERE id = $2
RETURNING *;

-- name: CheckPassword :one
SELECT id, username, created_at, updated_at, nickname
FROM users
WHERE password = $1;

-- name: FindUserByID :one
SELECT username, password, nickname
FROM users
WHERE id = $1;

-- name: FindUserByUsername :one
SELECT id, password, nickname
FROM users
WHERE username = $1;