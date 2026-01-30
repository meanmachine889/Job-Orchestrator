package store

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type JobDetail struct {
	ID             uuid.UUID
	Type           string
	Payload        json.RawMessage
	Status         string
	RetryCount     int
	MaxRetries     int
	TimeoutSeconds int
	WorkerID       *uuid.UUID
	Error          *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *Store) GetJobDetail(ctx context.Context, jobId uuid.UUID) (*JobDetail, error) {
	query := `SELECT * FROM jobs WHERE id = $1`
	row, err := s.db.QueryContext(ctx, query, jobId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, nil
	}

	var job JobDetail
	if err := row.Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.RetryCount,
		&job.MaxRetries,
		&job.TimeoutSeconds,
		&job.WorkerID,
		&job.Error,
		&job.CreatedAt,
		&job.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &job, nil
}
