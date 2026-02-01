package store

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type WorkerRow struct {
	ID            uuid.UUID
	Hostname      string
	Status        string
	LastHeartbeat time.Time
}

func (s *Store) ListWorkers(ctx context.Context) ([]*WorkerRow, error) {
	query := `SELECT id, hostname, status, last_heartbeat FROM workers ORDER BY last_heartbeat DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workers []*WorkerRow
	for rows.Next() {
		var w WorkerRow
		if err := rows.Scan(&w.ID, &w.Hostname, &w.Status, &w.LastHeartbeat); err != nil {
			return nil, err
		}
		workers = append(workers, &w)
	}
	return workers, nil
}
