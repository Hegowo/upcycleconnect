package main

import (
	"fmt"
	"log"
	"os"

	"upcycleconnect/backend/internal/config"
	"upcycleconnect/backend/internal/database"
	"upcycleconnect/backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := router.Setup(db, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
