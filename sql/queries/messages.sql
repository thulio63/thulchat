-- name: SendMessage :one
INSERT INTO messages (sender_id, body, sent_at, server_id)
VALUES (
    $1,
    $2,
    NOW(),
    $3
)
RETURNING *;

-- name: RetrieveMessages :many
SELECT sender_id, body, sent_at
FROM messages
WHERE server_id = $1
ORDER BY sent_at;