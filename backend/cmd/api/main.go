package main

import (
	"fmt"
	"log"
	"os"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/models"
	"pengajuan-dokumen/backend/internal/routes"
)

func main() {
	log.Println("Starting Pengajuan Dokumen Backend API...")

	// 1. Initialize Database
	db := config.ConnectDatabase()

	// 2. Auto Migrate Models
	log.Println("Running AutoMigrate database models...")
	err := db.AutoMigrate(
		&models.User{},
		&models.DocumentType{},
		&models.Project{},
		&models.Document{},
		&models.ReviewHistory{},
	)
	if err != nil {
		log.Fatalf("Database AutoMigrate failed: %v", err)
	}
	log.Println("Database migration completed successfully.")

	// 3. Create uploads directory for file storage
	if err := os.MkdirAll("uploads", 0750); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}

	// 4. Setup Gin Router
	r := routes.SetupRouter()

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://127.0.0.1:%s\n", port)
	if err := r.Run(fmt.Sprintf("127.0.0.1:%s", port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
