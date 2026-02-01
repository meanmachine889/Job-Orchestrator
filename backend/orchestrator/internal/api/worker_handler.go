package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type registerWorkerRequest struct {
	Hostname string `json:"hostname"`
}

type heartbeatRequest struct {
	ID string `json:"id"`
}

type WorkerDTO struct {
	ID            string    `json:"id"`
	Hostname      string    `json:"hostname"`
	Status        string    `json:"status"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

func (h *Handler) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	var req registerWorkerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	worker, err := h.store.CreateWorker(ctx, req.Hostname)
	if err != nil {
		http.Error(w, "Failed to create worker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":     worker.ID.String(),
		"status": worker.Status,
	})
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var req heartbeatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	workerID, err := uuid.Parse(req.ID)

	if err != nil {
		http.Error(w, "Invalid worker ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = h.store.UpdateHeartbeat(ctx, workerID)
	if err != nil {
		http.Error(w, "Failed to update heartbeat", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ONLINE",
	})
}

func (h *Handler) ListWorkers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	workers, err := h.store.ListWorkers(ctx)
	if err != nil {
		http.Error(w, "Failed to list workers", http.StatusInternalServerError)
		return
	}

	var workerDTOs []WorkerDTO
	for _, w := range workers {
		workerDTOs = append(workerDTOs, WorkerDTO{
			ID:            w.ID.String(),
			Hostname:      w.Hostname,
			Status:        w.Status,
			LastHeartbeat: w.LastHeartbeat,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"workers": workerDTOs,
	})
}
