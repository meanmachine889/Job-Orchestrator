package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

// defines how our job creation request looks like
type createJobRequest struct {
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	MaxRetries     int             `json:"max_retries"`
	TimeoutSeconds int             `json:"timeout_seconds"`
}

type JobDTO struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
	WorkerID   *string   `json:"worker_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type JobDetailDTO struct {
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	Status         string          `json:"status"`
	RetryCount     int             `json:"retry_count"`
	MaxRetries     int             `json:"max_retries"`
	TimeoutSeconds int             `json:"timeout_seconds"`
	WorkerID       *string         `json:"worker_id"`
	Error          *string         `json:"error"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	//made to handle job creation requests
	var req createJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job := &store.JobCreate{
		ID:             uuid.New(),
		Type:           req.Type,
		Payload:        req.Payload,
		Status:         "PENDING",
		MaxRetries:     req.MaxRetries,
		TimeoutSeconds: req.TimeoutSeconds,
	}

	// Creates a new context with a 3-second timeout derived from the HTTP request's context.
	// The cancel function should be deferred to release resources when the operation completes
	// or when the timeout expires, whichever occurs first.
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = h.store.CreateJob(ctx, job)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	err = h.queue.Enqueue(ctx, job.ID.String())
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

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit := 20
	offset := 0

	l := query.Get("limit")
	o := query.Get("offset")

	if l != "" {
		lim, err := strconv.Atoi(l)
		if err == nil && lim > 0 {
			limit = lim
		}
	}

	if o != "" {
		off, err := strconv.Atoi(o)
		if err == nil && off >= 0 {
			offset = off
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	jobs, err := h.store.ListJobs(ctx, limit, offset)

	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	var response []JobDTO

	for _, job := range jobs {
		dto := JobDTO{
			ID:         job.ID.String(),
			Type:       job.Type,
			Status:     job.Status,
			RetryCount: job.RetryCount,
			CreatedAt:  job.CreatedAt,
		}
		response = append(response, dto)
	}

	resp := map[string]interface{}{
		"jobs": response,
		"limit": limit,
		"offset": offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h* Handler) GetJobDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/jobs/")
	jobId, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	job, err := h.store.GetJobDetail(ctx, jobId)
	if err != nil {
		http.Error(w, "Failed to fetch job detail", http.StatusInternalServerError)
		return
	}
	if job == nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	var workerID *string
	if job.WorkerID != nil {
		wid := job.WorkerID.String()
		workerID = &wid
	}

	dto := JobDetailDTO{
		ID:             job.ID.String(),
		Type:           job.Type,
		Payload:        job.Payload,
		Status:         job.Status,
		RetryCount:     job.RetryCount,
		MaxRetries:     job.MaxRetries,
		TimeoutSeconds: job.TimeoutSeconds,
		WorkerID:       workerID,
		Error:          job.Error,
		CreatedAt:      job.CreatedAt,
		UpdatedAt:      job.UpdatedAt,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}