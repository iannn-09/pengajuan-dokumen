package main

import (
	"log"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/seeder"
)

func main() {
	log.Println("Starting Standalone Database Seeder Tool...")

	// Connect Database
	config.ConnectDatabase()

	// Run Seeder: 1,000 Pemohon + 1,000 Penilai + 10,000 Projects
	seeder.SeedDatabase(1000, 1000, 10000)

	log.Println("Seeder CLI finished successfully!")
}
