package store

import (
	"context"
	"time"
	
	"github.com/google/uuid"
)

type JobRow struct {
	ID         uuid.UUID
	Type       string
	Status     string
	RetryCount int
	WorkerID   *uuid.UUID
	CreatedAt  time.Time
}

func (s *Store) ListJobs(
	ctx context.Context,
	limit int,
	offset int,
) ([]JobRow, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id,
			type,
			status,
			retry_count,
			worker_id,
			created_at
		FROM jobs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []JobRow
	for rows.Next() {
		var job JobRow
		if err := rows.Scan(
			&job.ID,
			&job.Type,
			&job.Status,
			&job.RetryCount,
			&job.WorkerID,
			&job.CreatedAt,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

