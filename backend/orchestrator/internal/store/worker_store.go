package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	ID            uuid.UUID
	Hostname      string
	Status        string
	LastHeartbeat time.Time
}

func (s *Store) CreateWorker(ctx context.Context, hostname string) (*Worker, error) {
	worker := &Worker{
		ID:            uuid.New(),
		Hostname:      hostname,
		Status:        "ONLINE",
		LastHeartbeat: time.Now(),
	}

	query := `INSERT INTO workers (id, hostname, status, last_heartbeat) VALUES ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, worker.ID, worker.Hostname, worker.Status, worker.LastHeartbeat)
	if err != nil {
		return nil, err
	}
	return worker, nil
}

func (s *Store) UpdateHeartbeat(ctx context.Context, workerID uuid.UUID) error {
	query := `UPDATE workers SET last_heartbeat = $1, status = $2 WHERE id = $3`
	_, err := s.db.ExecContext(ctx, query, time.Now(), "ONLINE", workerID)
	return err
}

func (s *Store) MarkWorkerOffline(ctx context.Context, timeout time.Duration) error {
	query := `UPDATE workers SET status = $1 WHERE last_heartbeat < NOW() - $2::interval`
	_, err := s.db.ExecContext(ctx, query, "OFFLINE", timeout.String())
	return err
}
