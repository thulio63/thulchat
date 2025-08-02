-- name: SendMessage :one
INSERT INTO messages (sender_id, body, sent_at, hostname, port)
VALUES (
    $1,
    $2,
    NOW(),
    $3, 
    $4
)
RETURNING *;