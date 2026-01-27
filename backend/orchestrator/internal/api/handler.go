package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/queue"
	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/store"
)

type Handler struct {
	store *store.Store
	queue *queue.Queue
}

func NewHandler(store *store.Store, queue *queue.Queue) *Handler {
	return &Handler{
		store: store,
		queue: queue,
	}
}

type createJobRequest struct {
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	MaxRetries     int             `json:"max_retries"`
	TimeoutSeconds int             `json:"timeout_seconds"`
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job := &store.Job{
		ID:             uuid.New(),
		Type:           req.Type,
		Payload:        req.Payload,
		Status:         "PENDING",
		MaxRetries:     req.MaxRetries,
		TimeoutSeconds: req.TimeoutSeconds,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = h.store.CreateJob(ctx, job)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	err = h.queue.Enqueue(ctx, "job_queue", job.ID.String())
	if err != nil {
		http.Error(w, "Failed to enqueue job", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"id":     job.ID.String(),
		"status": job.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
