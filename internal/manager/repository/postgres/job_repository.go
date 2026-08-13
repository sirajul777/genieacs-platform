package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/sirajul777/genieacs-platform/internal/manager/domain/job"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository/postgres/sqlc"
)

var ErrNoPendingJob = repository.ErrNoPendingJob

type JobRepository struct { pool *pgxpool.Pool; queries *sqlc.Queries }
func NewJobRepository(pool *pgxpool.Pool) *JobRepository { return &JobRepository{pool: pool, queries: sqlc.New(pool)} }
func (r *JobRepository) Create(ctx context.Context, item domain.Job) (domain.Job, error) {
	created, err := r.queries.CreateJob(ctx, sqlc.CreateJobParams{ID: item.ID, Type: item.Type, Payload: item.Payload})
	if err != nil { return domain.Job{}, err }
	return jobToDomain(created), nil
}
func (r *JobRepository) AssignNext(ctx context.Context, agentID string) (domain.Job, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return domain.Job{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	row := tx.QueryRow(ctx, `UPDATE jobs SET status = 'in_progress', agent_id = $1, updated_at = NOW() WHERE id = (SELECT id FROM jobs WHERE status = 'pending' ORDER BY created_at ASC FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING id, type, payload, status, agent_id, progress, result, error, created_at, updated_at`, agentID)
	var item sqlc.Job
	if err := row.Scan(&item.ID, &item.Type, &item.Payload, &item.Status, &item.AgentID, &item.Progress, &item.Result, &item.Error, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return domain.Job{}, repository.ErrNoPendingJob }
		return domain.Job{}, err
	}
	if err := tx.Commit(ctx); err != nil { return domain.Job{}, err }
	return jobToDomain(item), nil
}
func (r *JobRepository) Progress(ctx context.Context, id string, progress int) (domain.Job, error) {
	item, err := r.queries.UpdateJobProgress(ctx, id, progress)
	if err != nil { return domain.Job{}, err }
	return jobToDomain(item), nil
}
func (r *JobRepository) Complete(ctx context.Context, id string, result []byte) (domain.Job, error) {
	item, err := r.queries.CompleteJob(ctx, id, result)
	if err != nil { return domain.Job{}, err }
	return jobToDomain(item), nil
}
func (r *JobRepository) Fail(ctx context.Context, id string, message string) (domain.Job, error) {
	item, err := r.queries.FailJob(ctx, id, message)
	if err != nil { return domain.Job{}, err }
	return jobToDomain(item), nil
}
func jobToDomain(item sqlc.Job) domain.Job { return domain.Job{ID:item.ID, Type:item.Type, Payload:item.Payload, Status:domain.Status(item.Status), AgentID:item.AgentID, Progress:item.Progress, Result:item.Result, Error:item.Error, CreatedAt:item.CreatedAt, UpdatedAt:item.UpdatedAt} }

type JobEventRepository struct{ queries *sqlc.Queries }
func NewJobEventRepository(pool *pgxpool.Pool) *JobEventRepository { return &JobEventRepository{queries: sqlc.New(pool)} }
func (r *JobEventRepository) Create(ctx context.Context, jobID, eventType, message string) (domain.Event, error) {
	item, err := r.queries.CreateJobEvent(ctx, jobID, eventType, message)
	if err != nil { return domain.Event{}, err }
	return domain.Event{ID:item.ID, JobID:item.JobID, Type:item.Type, Message:item.Message, CreatedAt:item.CreatedAt}, nil
}
