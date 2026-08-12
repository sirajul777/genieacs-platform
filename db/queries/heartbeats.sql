-- name: CreateHeartbeat :one
INSERT INTO heartbeats (agent_id)
VALUES ($1)
RETURNING id, agent_id, created_at;
