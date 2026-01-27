package job

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Pending  Status = "PENDING"
	Running  Status = "RUNNING"
	Success  Status = "SUCCESS"
	Failed   Status = "FAILED"
	Retrying Status = "RETRYING"
	Dead     Status = "DEAD"
)

type Job struct {
	ID      uuid.UUID       `db:"id"`
	Type    string          `db:"type"`
	Payload json.RawMessage `db:"payload"`
	Status  Status          `db:"status"`

	RetryCount     int `db:"retry_count"`
	MaxRetries     int `db:"max_retries"`
	TimeoutSeconds int `db:"timeout_seconds"`

	WorkerID *uuid.UUID `db:"worker_id"`
	Error    *string    `db:"error"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
