package main

import (
	"log"

	"ticket-box-be/internal/admin"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database for migration: %v", err)
	}

	// 1. Migrate Domain Models (tickets, etc.)
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// 2. Migrate GoAdmin Internal Schema
	if err := admin.InitSchema(db, cfg); err != nil {
		log.Fatalf("❌ GoAdmin schema initialization failed: %v", err)
	}

	log.Println("🎉 Migration finished successfully!")
}
