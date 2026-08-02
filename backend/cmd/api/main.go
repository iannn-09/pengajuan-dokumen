package main

import (
	"log"
	"os"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/models"
	"pengajuan-dokumen/backend/internal/routes"
	"pengajuan-dokumen/backend/internal/services"
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

	// 3. Initialize WhatsApp Native Whatsmeow Service
	services.InitWhatsAppService()

	// 4. Create uploads directory for file storage
	if err := os.MkdirAll("uploads", 0750); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}

	// 5. Setup Gin Router
	r := routes.SetupRouter()

	// 6. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
