-- name: CreateAgent :one
INSERT INTO agents (id, name, endpoint, status, token_hash)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, endpoint, status, token_hash, created_at, updated_at, last_seen;

-- name: GetAgentByID :one
SELECT id, name, endpoint, status, token_hash, created_at, updated_at, last_seen
FROM agents
WHERE id = $1;

-- name: ListAgents :many
SELECT id, name, endpoint, status, token_hash, created_at, updated_at, last_seen
FROM agents
ORDER BY created_at DESC;

-- name: UpdateAgentLastSeen :one
UPDATE agents
SET last_seen = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING id, name, endpoint, status, token_hash, created_at, updated_at, last_seen;
