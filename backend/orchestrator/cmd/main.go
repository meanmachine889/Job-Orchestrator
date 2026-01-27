package main

import "log"
import "github.com/meanmachine889/distributed-orchestrator/orchestrator/internal/store"

func main() {
	_, err := store.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected")

	// pass db to API, scheduler, etc
}