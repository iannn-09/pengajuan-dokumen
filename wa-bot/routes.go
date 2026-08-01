package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// SetupRoutes sets up all API routes
func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Status endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sessions := sessionManager.ListSessions()
		connectedCount := 0
		for _, s := range sessions {
			if s.Status == "connected" {
				connectedCount++
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":             "Multi-Session WhatsApp Bot API is running",
			"total_sessions":     len(sessions),
			"connected_sessions": connectedCount,
			"timestamp":          time.Now().Format(time.RFC3339),
		})
	}).Methods("GET")

	// Session management endpoints
	r.HandleFunc("/sessions", handleCreateSession).Methods("POST")
	r.HandleFunc("/sessions", handleListSessions).Methods("GET")
	r.HandleFunc("/sessions/{id}", handleGetSession).Methods("GET")
	r.HandleFunc("/sessions/{id}", handleUpdateSession).Methods("PATCH")
	r.HandleFunc("/sessions/{id}/logout", handleLogoutSession).Methods("POST")
	r.HandleFunc("/sessions/{id}", handleDeleteSession).Methods("DELETE")

	// Message endpoints
	r.HandleFunc("/send-message", handleSendMessage).Methods("POST")
	r.HandleFunc("/send-bulk-same-message", handleBulkSendSameMessage).Methods("POST")
	r.HandleFunc("/send-bulk-different-messages", handleBulkSendDifferentMessages).Methods("POST")
	r.HandleFunc("/groups", handleGetGroups).Methods("GET")

	// CORS middleware
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(r)

	return handler
}
