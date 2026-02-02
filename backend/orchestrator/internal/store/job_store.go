package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type JobCreate struct {
	ID             uuid.UUID
	Type           string
	Payload        json.RawMessage
	Status         string
	RetryCount     int
	MaxRetries     int
	TimeoutSeconds int
}

func (s *Store) CreateJob(ctx context.Context, job *JobCreate) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO jobs (
id, type, payload, status, max_retries, timeout_seconds
) VALUES ($1, $2, $3, $4, $5, $6)`,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.MaxRetries,
		job.TimeoutSeconds,
	)
	return err
}

func (s *Store) AssignNextJob(ctx context.Context, workerID uuid.UUID) (*JobCreate, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var job JobCreate

	// query selects the oldest pending job from the jobs table for processing.
	// It uses row-level locking (FOR UPDATE) to prevent concurrent access,
	// and SKIP LOCKED to avoid waiting on already-locked rows, enabling
	// multiple workers to efficiently pick up different pending jobs simultaneously.
	query := `SELECT id, type, payload, retry_count, max_retries FROM jobs WHERE status = 'PENDING' ORDER BY created_at LIMIT 1 FOR UPDATE SKIP LOCKED`

	err = tx.QueryRowContext(ctx, query).Scan(&job.ID, &job.Type, &job.Payload, &job.RetryCount, &job.MaxRetries)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE jobs SET status = 'RUNNING', worker_id = $1, updated_at = NOW() WHERE id = $2`,
		workerID,
		job.ID,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *Store) ReportJobResult(ctx context.Context, jobID uuid.UUID, status string, errMsg string) error {
	query := `UPDATE jobs SET status = $1, error = $2, updated_at = NOW() WHERE id = $3`
	_, err := s.db.ExecContext(ctx, query, status, errMsg, jobID)
	return err
}

func (s *Store) HandleJobFailures(ctx context.Context, jobId uuid.UUID, errormsg string) (error, bool) {
	var retrycount, max_retries int

	err := s.db.QueryRowContext(ctx,
		`SELECT retry_count, max_retries FROM jobs WHERE id = $1`,
		jobId,
	).Scan(&retrycount, &max_retries)

	if err != nil {
		return err, false
	}

	if retrycount+1 > max_retries {
		_, err = s.db.ExecContext(ctx,
			`UPDATE jobs SET status = 'DEAD', error = $1, updated_at = NOW() WHERE id = $2`,
			errormsg,
			jobId,
		)
		return err, false
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE jobs SET status = 'PENDING', retry_count = retry_count + 1, error = $1, updated_at = NOW() WHERE id = $2`,
		errormsg,
		jobId,
	)
	if err != nil {
		return err, false
	}

	return nil, true
}
