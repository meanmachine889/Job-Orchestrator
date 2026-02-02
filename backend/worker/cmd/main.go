package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/meanmachine889/distributed-orchestrator/worker/internal/executor"
	"github.com/meanmachine889/distributed-orchestrator/worker/internal/orchestrator"
	redisclient "github.com/meanmachine889/distributed-orchestrator/worker/internal/redis"
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

	redisClient := redisclient.New(os.Getenv("REDIS_URL"))

	go func() {
		for {

			jobId, err := redisClient.WaitForJob(context.Background())

			if err != nil {
				log.Println("Redis wait failed:", err)
				continue
			}

			log.Println("Woken by Redis for job:", jobId)
			job, err := client.FetchJob(workerId)
			if err != nil {
				log.Println("Failed to fetch job:", err)
				continue
			}
			if job == nil {
				continue
			}
			attempt := job.RetryCount + 1
			if job.RetryCount > 0 {
				log.Printf("Retrying job %s (attempt %d/%d)", job.ID, attempt, job.MaxRetries+1)
			} else {
				log.Printf("Executing job: %s (max retries: %d)", job.ID, job.MaxRetries)
			}
			err = executor.Execute(job.Type)
			if err != nil {
				if attempt == job.MaxRetries+1 {
					log.Printf("Job %s DEAD after %d attempts: %v", job.ID, attempt, err)
				} else {
					log.Printf("Job %s FAILED on attempt %d/%d: %v", job.ID, attempt, job.MaxRetries+1, err)
				}
				time.Sleep(time.Duration(2^attempt) * 2*time.Second)
				client.ReportJobResult(job.ID.String(), "FAILED", err.Error())
			} else {
				client.ReportJobResult(job.ID.String(), "SUCCESS", "")
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			err := client.SendHeartbeat(workerId)
			if err != nil {
				log.Println("Failed to send heartbeat:", err)
			}
		case <-sig:
			log.Println("Shutting down worker...")
			return
		}
	}
}
