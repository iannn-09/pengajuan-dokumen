package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var globalLogger waLog.Logger

func syncLatestWAVersion() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	latestVer, err := whatsmeow.GetLatestVersion(ctx, nil)
	if err != nil {
		log.Printf("Failed to fetch latest WhatsApp Web version, using bundled version: %s (err: %v)", store.GetWAVersion().String(), err)
		return
	}

	currentVer := store.GetWAVersion()
	if currentVer.LessThan(*latestVer) {
		store.SetWAVersion(*latestVer)
		log.Printf("WhatsApp Web version updated at runtime: %s -> %s", currentVer.String(), latestVer.String())
	} else {
		log.Printf("WhatsApp Web version already up to date: %s", currentVer.String())
	}
}

func main() {
	globalLogger = waLog.Stdout("whatsapp", "INFO", true)

	syncLatestWAVersion()

	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	} else {
		log.Println("Loaded .env file")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL != "" {
		log.Printf("Webhook URL configured: %s", webhookURL)
	}

	// Create session directory
	if err := os.MkdirAll("session", 0755); err != nil {
		log.Fatalf("Failed to create session directory: %v", err)
	}

	// Initialize session manager
	var err error
	sessionManager, err = NewSessionManager(globalLogger, webhookURL)
	if err != nil {
		log.Fatalf("Failed to initialize session manager: %v", err)
	}

	// Load persisted sessions
	if err := sessionManager.LoadPersistedSessions(); err != nil {
		log.Printf("Warning: Failed to load persisted sessions: %v", err)
	}

	// Setup routes
	handler := SetupRoutes()

	log.Println("Multi-Session WhatsApp Bot API starting on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
