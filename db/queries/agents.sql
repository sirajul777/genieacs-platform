-- name: CreateAgent :one
INSERT INTO agents (id, name, endpoint, status)
VALUES ($1, $2, $3, $4)
RETURNING id, name, endpoint, status, created_at, updated_at;

-- name: GetAgentByID :one
SELECT id, name, endpoint, status, created_at, updated_at
FROM agents
WHERE id = $1;

-- name: ListAgents :many
SELECT id, name, endpoint, status, created_at, updated_at
FROM agents
ORDER BY created_at DESC;
