package store

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Job struct {
	ID             uuid.UUID
	Type           string
	Payload        json.RawMessage
	Status         string
	MaxRetries     int
	TimeoutSeconds int
}

func (s *Store) CreateJob(ctx context.Context, job *Job) error {
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
