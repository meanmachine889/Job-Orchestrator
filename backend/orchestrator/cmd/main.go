package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/api"
	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/queue"
	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/scheduler"
	"github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/store"
)

func main() {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()
	db, err := store.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	jobQueue := queue.New(os.Getenv("REDIS_URL"))
	handler := api.NewHandler(db, jobQueue)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	corsHandler := api.CorsMiddleware(mux)

	monitor := scheduler.NewWorkerMonitor(db)
	go monitor.Start()

	log.Println("Orchestrator listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
