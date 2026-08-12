-- name: CreateJobEvent :one
INSERT INTO job_events (job_id, type, message)
VALUES ($1, $2, $3)
RETURNING id, job_id, type, message, created_at;
