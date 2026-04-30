package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/database"
	"upcycleconnect/backend/models"
	"upcycleconnect/backend/router"
)

func main() {
	cfg := config.Load()

	log.Printf("[config] APP_URL=%s", cfg.AppURL)
	if strings.Contains(cfg.AppURL, "localhost") || strings.Contains(cfg.AppURL, "127.0.0.1") {
		log.Printf("[config] WARNING: APP_URL points to localhost — Stripe redirects after payment will fail in production. Set APP_URL=https://upcycleconnect.xyz in .env on the prod server.")
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Reservation{},
		&models.Invoice{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
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
