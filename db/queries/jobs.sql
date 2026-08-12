-- name: CreateJob :one
INSERT INTO jobs (id, type, payload, status)
VALUES ($1, $2, $3, 'pending')
RETURNING id, type, payload, status, agent_id, progress, result, error, created_at, updated_at;

-- name: AssignNextJob :one
UPDATE jobs
SET status = 'in_progress', agent_id = $1, updated_at = NOW()
WHERE id = (
  SELECT id FROM jobs
  WHERE status = 'pending'
  ORDER BY created_at ASC
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
RETURNING id, type, payload, status, agent_id, progress, result, error, created_at, updated_at;

-- name: UpdateJobProgress :one
UPDATE jobs
SET progress = $2, updated_at = NOW()
WHERE id = $1 AND status = 'in_progress'
RETURNING id, type, payload, status, agent_id, progress, result, error, created_at, updated_at;

-- name: CompleteJob :one
UPDATE jobs
SET status = 'completed', progress = 100, result = $2, updated_at = NOW()
WHERE id = $1 AND status = 'in_progress'
RETURNING id, type, payload, status, agent_id, progress, result, error, created_at, updated_at;

-- name: FailJob :one
UPDATE jobs
SET status = 'failed', error = $2, updated_at = NOW()
WHERE id = $1 AND status = 'in_progress'
RETURNING id, type, payload, status, agent_id, progress, result, error, created_at, updated_at;
