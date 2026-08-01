package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

var sessionManager *SessionManager

// Session Management Handlers

func handleCreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	if req.Name == "" {
		req.Name = "Session " + time.Now().Format("2006-01-02 15:04:05")
	}

	session, err := sessionManager.CreateSession(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Wait a bit for QR code generation
	time.Sleep(1 * time.Second)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      session.ID,
		"name":    session.Name,
		"qr_code": session.QRCode,
		"status":  session.Status,
	})
}

func handleListSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sessions := sessionManager.ListSessions()
	result := make([]map[string]interface{}, 0, len(sessions))

	for _, session := range sessions {
		session.mu.RLock()
		result = append(result, map[string]interface{}{
			"id":                session.ID,
			"name":              session.Name,
			"phone_number":      session.PhoneNumber,
			"status":            session.Status,
			"last_connected_at": session.LastConnected,
		})
		session.mu.RUnlock()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions": result,
	})
}

func handleGetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, ok := sessionManager.GetSession(sessionID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Session not found"})
		return
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                session.ID,
		"name":              session.Name,
		"phone_number":      session.PhoneNumber,
		"status":            session.Status,
		"qr_code":           session.QRCode,
		"last_connected_at": session.LastConnected,
	})
}

func handleUpdateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sessionID := vars["id"]

	var req updateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	if err := sessionManager.UpdateSessionName(sessionID, req.Name); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session updated"})
}

func handleLogoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sessionID := vars["id"]

	var req sessionSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	if err := sessionManager.LogoutSession(sessionID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session logged out"})
}

func handleDeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sessionID := vars["id"]

	var req sessionSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	if err := sessionManager.DeleteSession(sessionID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session deleted"})
}

// Message Handlers

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Use first available connected session
	sessions := sessionManager.ListSessions()
	var client *whatsmeow.Client
	var sessionID string

	for _, session := range sessions {
		if session.Client.IsConnected() && session.Client.Store.ID != nil {
			client = session.Client
			sessionID = session.ID
			break
		}
	}

	if client == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "No connected WhatsApp session available"})
		return
	}

	sendSuratTugas := r.URL.Query().Get("surattugas") == "true"

	var targetJID types.JID
	var targetType string

	if strings.Contains(req.Target, "@g.us") {
		var err error
		targetJID, err = types.ParseJID(req.Target)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid group JID format: " + err.Error()})
			return
		}
		targetType = "group"
	} else {
		normalizedTarget := normalizePhoneNumber(req.Target)
		targetJID = types.NewJID(normalizedTarget, types.DefaultUserServer)
		targetType = "personal"
	}

	var err error
	if sendSuratTugas {
		err = sendSuratTugasPDFWithCaption(client, targetJID, req.Message)
	} else {
		err = sendMessageWithRetry(client, targetJID, req.Message, 3)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "Success",
		"session_id":  sessionID,
		"target":      req.Target,
		"target_type": targetType,
	})
}

func handleBulkSendSameMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req bulkMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Get connected sessions
	var sessions []*WhatsAppSession
	if len(req.SessionIDs) > 0 {
		sessions = sessionManager.GetConnectedSessions(req.SessionIDs)
	} else {
		allSessions := sessionManager.ListSessions()
		for _, s := range allSessions {
			if s.Client.IsConnected() && s.Client.Store.ID != nil {
				sessions = append(sessions, s)
			}
		}
	}

	if len(sessions) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "No connected sessions available"})
		return
	}

	strategy := req.Strategy
	if strategy == "" {
		strategy = "single"
	}

	results := make([]map[string]interface{}, len(req.Targets))
	sessionIndex := 0

	for i, target := range req.Targets {
		// Select session based on strategy
		var selectedSession *WhatsAppSession
		switch strategy {
		case "even":
			selectedSession = sessions[sessionIndex%len(sessions)]
			sessionIndex++
		case "random":
			idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(sessions))))
			selectedSession = sessions[idx.Int64()]
		default: // "single"
			selectedSession = sessions[0]
		}

		var jid types.JID
		var targetType string

		if strings.Contains(target, "@g.us") {
			var err error
			jid, err = types.ParseJID(target)
			if err != nil {
				results[i] = map[string]interface{}{
					"target":     target,
					"success":    false,
					"error":      "Invalid group JID format: " + err.Error(),
					"session_id": selectedSession.ID,
				}
				continue
			}
			targetType = "group"
		} else {
			normalizedTarget := normalizePhoneNumber(target)
			jid = types.NewJID(normalizedTarget, types.DefaultUserServer)
			targetType = "personal"
		}

		var err error
		if req.SendSuratTugas {
			// Send with PDF attachment
			err = sendSuratTugasPDFWithCaption(selectedSession.Client, jid, req.Message)
		} else {
			// Send regular text message
			err = sendMessageWithRetry(selectedSession.Client, jid, req.Message, 2)
		}

		results[i] = map[string]interface{}{
			"target":           target,
			"target_jid":       jid.String(),
			"target_type":      targetType,
			"success":          err == nil,
			"session_id":       selectedSession.ID,
			"send_surat_tugas": req.SendSuratTugas,
		}

		if err != nil {
			results[i]["error"] = err.Error()
		}

		if i < len(req.Targets)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "Bulk same message processing completed",
		"results":  results,
		"strategy": strategy,
	})
}

func handleBulkSendDifferentMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req bulkDifferentMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if SECRET == "" {
		SECRET = "default-secret"
	}

	if req.Secret != SECRET {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	results := make([]map[string]interface{}, len(req.Messages))

	for i, msg := range req.Messages {
		// Get the assigned session from BE
		session, exists := sessionManager.GetSession(msg.SessionID)
		if !exists || session == nil || !session.Client.IsConnected() {
			results[i] = map[string]interface{}{
				"target":     msg.Targets,
				"success":    false,
				"error":      "Session not found or not connected: " + msg.SessionID,
				"message":    msg.Message,
				"session_id": msg.SessionID,
			}
			continue
		}

		var jid types.JID
		var targetType string

		if strings.Contains(msg.Targets, "@g.us") {
			var err error
			jid, err = types.ParseJID(msg.Targets)
			if err != nil {
				results[i] = map[string]interface{}{
					"target":     msg.Targets,
					"success":    false,
					"error":      "Invalid group JID format: " + err.Error(),
					"message":    msg.Message,
					"session_id": session.ID,
				}
				continue
			}
			targetType = "group"
		} else {
			normalizedTarget := normalizePhoneNumber(msg.Targets)
			jid = types.NewJID(normalizedTarget, types.DefaultUserServer)
			targetType = "personal"
		}

		var err error
		if req.SendSuratTugas {
			// Send with PDF attachment
			err = sendSuratTugasPDFWithCaption(session.Client, jid, msg.Message)
		} else {
			// Send regular text message
			err = sendMessageWithRetry(session.Client, jid, msg.Message, 2)
		}

		results[i] = map[string]interface{}{
			"target":           msg.Targets,
			"target_jid":       jid.String(),
			"target_type":      targetType,
			"success":          err == nil,
			"message":          msg.Message,
			"session_id":       session.ID,
			"send_surat_tugas": req.SendSuratTugas,
		}

		if err != nil {
			results[i]["error"] = err.Error()
		}

		if i < len(req.Messages)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "Bulk different messages processing completed",
		"results": results,
	})
}

func handleGetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Use first connected session
	sessions := sessionManager.ListSessions()
	var client *whatsmeow.Client

	for _, session := range sessions {
		if session.Client.IsConnected() && session.Client.Store.ID != nil {
			client = session.Client
			break
		}
	}

	if client == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "No connected session available"})
		return
	}

	joinedGroups, err := client.GetJoinedGroups(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve groups: " + err.Error()})
		return
	}

	var groups []map[string]interface{}
	for _, groupInfo := range joinedGroups {
		group := map[string]interface{}{
			"gid":  groupInfo.JID.String(),
			"name": groupInfo.GroupName.Name,
		}
		groups = append(groups, group)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "Success",
		"total_groups": len(groups),
		"groups":       groups,
	})
}
