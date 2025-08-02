-- name: CreateServer :one
INSERT INTO servers (creator_id, server_id, created_at, hostname, port)
VALUES (
    $1,
    gen_random_uuid(),
    NOW(),
    $2, 
    $3
)
RETURNING *;

-- name: DeleteServer :one
DELETE FROM servers
WHERE creator_id = $1
RETURNING *;

-- name: RetrieveServers :many
SELECT creator_id, server_id, hostname, port 
FROM servers
ORDER BY server_id;