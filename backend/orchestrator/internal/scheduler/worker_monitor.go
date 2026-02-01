package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/store"
)

type WorkerMonitor struct {
	store *store.Store
}

func NewWorkerMonitor(store *store.Store) *WorkerMonitor {
	return &WorkerMonitor{store: store}
}

func (wm *WorkerMonitor) Start() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := wm.store.MarkWorkerOffline(ctx, 15*time.Second)
		if err != nil {
			log.Println("Failed to mark offline workers:", err)
		}
		cancel()
	}
}