package main

import (
	"log"

	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database for migration: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("🎉 Migration finished successfully!")
}
