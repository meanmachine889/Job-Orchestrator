package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/meanmachine889/distributed-orchestrator/worker/internal/orchestrator"
)

func main() {
	// Load .env from project root
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	orurl := os.Getenv("ORCHESTRATOR_URL")
	if orurl == "" {
		log.Fatal("ORCHESTRATOR_URL is not set")
	}
	client := orchestrator.New(orurl)
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	workerId, err := client.RegisterWorker(hostname)
	if err != nil {
		log.Fatal("Failed to register worker:", err)
	}

	log.Println("Registered worker with ID:", workerId)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			err := client.SendHeartbeat(workerId)
			if err != nil {
				log.Println("Failed to send heartbeat:", err)
			} else {
				log.Println("Heartbeat sent")
			}
		case <-sig:
			log.Println("Shutting down worker...")
			return
		}
	}
}
