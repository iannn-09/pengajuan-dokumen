package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/glebarez/sqlite"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// SessionManager manages multiple WhatsApp sessions
type SessionManager struct {
	sessions   map[string]*WhatsAppSession
	db         *sql.DB
	logger     waLog.Logger
	webhookURL string
	mu         sync.RWMutex
}

// Initialize session manager database
func initSessionManagerDB() (*sql.DB, error) {
	// Ensure session directory exists
	if err := os.MkdirAll("session", 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %v", err)
	}

	dbPath := filepath.Join("session", "sessions_manager.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open session manager database: %v", err)
	}

	// Create sessions table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS wa_sessions (
		id VARCHAR(6) PRIMARY KEY,
		name TEXT NOT NULL,
		phone_number TEXT,
		status TEXT NOT NULL,
		qr_code TEXT,
		created_at DATETIME NOT NULL,
		last_connected_at DATETIME
	);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create sessions table: %v", err)
	}

	return db, nil
}

// Generate 6-character alphanumeric session ID
func generateSessionID(db *sql.DB) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const idLength = 6

	for i := 0; i < 100; i++ { // Max 100 attempts to avoid infinite loop
		id := make([]byte, idLength)
		for j := 0; j < idLength; j++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
			if err != nil {
				return "", err
			}
			id[j] = chars[num.Int64()]
		}

		sessionID := string(id)

		// Check if ID already exists
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM wa_sessions WHERE id = ?", sessionID).Scan(&count)
		if err != nil {
			return "", err
		}

		if count == 0 {
			return sessionID, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique session ID after 100 attempts")
}

// NewSessionManager creates new session manager
func NewSessionManager(logger waLog.Logger, webhookURL string) (*SessionManager, error) {
	db, err := initSessionManagerDB()
	if err != nil {
		return nil, err
	}

	sm := &SessionManager{
		sessions:   make(map[string]*WhatsAppSession),
		db:         db,
		logger:     logger,
		webhookURL: webhookURL,
	}

	return sm, nil
}

// CreateSession creates new WhatsApp session
func (sm *SessionManager) CreateSession(name string) (*WhatsAppSession, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Generate unique session ID
	sessionID, err := generateSessionID(sm.db)
	if err != nil {
		return nil, err
	}

	// Create session directory
	sessionDir := filepath.Join("session", sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %v", err)
	}

	// Initialize whatsmeow store
	dbPath := filepath.Join(sessionDir, "store.db")
	dsn := "file:" + dbPath + "?_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"

	container, err := sqlstore.New(context.Background(), "sqlite", dsn, sm.logger)
	if err != nil {
		os.RemoveAll(sessionDir)
		return nil, fmt.Errorf("failed to create whatsmeow store: %v", err)
	}

	// Get device store
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		container.Close()
		os.RemoveAll(sessionDir)
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	// Create WhatsApp client
	client := whatsmeow.NewClient(deviceStore, sm.logger)

	// Create session object
	session := &WhatsAppSession{
		ID:             sessionID,
		Name:           name,
		Status:         "qr_pending",
		Client:         client,
		Container:      container,
		CreatedAt:      time.Now(),
		IsGeneratingQR: false,
	}

	// Add event handlers
	client.AddEventHandler(sm.createEventHandler(session))
	client.AddEventHandler(sm.createReceiptHandler(session))

	// Save to database
	_, err = sm.db.Exec(`
		INSERT INTO wa_sessions (id, name, phone_number, status, qr_code, created_at, last_connected_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, sessionID, name, "", "qr_pending", "", session.CreatedAt, nil)

	if err != nil {
		container.Close()
		os.RemoveAll(sessionDir)
		return nil, fmt.Errorf("failed to save session to database: %v", err)
	}

	// Store in memory
	sm.sessions[sessionID] = session

	// Generate QR code
	go sm.generateQRCode(session)

	log.Printf("[%s] Session created: %s", sessionID, name)
	return session, nil
}

// LoadPersistedSessions loads persisted sessions on startup
func (sm *SessionManager) LoadPersistedSessions() error {
	rows, err := sm.db.Query("SELECT id, name, phone_number, status, created_at FROM wa_sessions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, phoneNumber, status string
		var createdAt time.Time
		var lastConnected *time.Time

		if err := rows.Scan(&id, &name, &phoneNumber, &status, &createdAt); err != nil {
			log.Printf("Failed to scan session row: %v", err)
			continue
		}

		// Check if session directory exists
		sessionDir := filepath.Join("session", id)
		dbPath := filepath.Join(sessionDir, "store.db")

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			log.Printf("[%s] Session store not found, skipping", id)
			continue
		}

		// Load session
		dsn := "file:" + dbPath + "?_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
		container, err := sqlstore.New(context.Background(), "sqlite", dsn, sm.logger)
		if err != nil {
			log.Printf("[%s] Failed to load store: %v", id, err)
			continue
		}

		deviceStore, err := container.GetFirstDevice(context.Background())
		if err != nil {
			log.Printf("[%s] Failed to get device: %v", id, err)
			container.Close()
			continue
		}

		client := whatsmeow.NewClient(deviceStore, sm.logger)

		session := &WhatsAppSession{
			ID:             id,
			Name:           name,
			PhoneNumber:    phoneNumber,
			Status:         "disconnected",
			Client:         client,
			Container:      container,
			CreatedAt:      createdAt,
			LastConnected:  lastConnected,
			IsGeneratingQR: false,
		}

		// Add event handlers
		client.AddEventHandler(sm.createEventHandler(session))
		client.AddEventHandler(sm.createReceiptHandler(session))

		// Store in memory
		sm.mu.Lock()
		sm.sessions[id] = session
		sm.mu.Unlock()

		// Try to reconnect
		go func(s *WhatsAppSession) {
			if s.Client.Store.ID != nil {
				log.Printf("[%s] Attempting to reconnect...", s.ID)
				if err := s.Client.Connect(); err != nil {
					log.Printf("[%s] Reconnection failed: %v", s.ID, err)
					sm.generateQRCode(s)
				} else {
					log.Printf("[%s] Reconnected successfully", s.ID)
				}
			} else {
				log.Printf("[%s] No existing session, generating QR", s.ID)
				sm.generateQRCode(s)
			}
		}(session)
	}

	return nil
}

// generateQRCode generates QR code for session
func (sm *SessionManager) generateQRCode(session *WhatsAppSession) {
	session.mu.Lock()
	if session.IsGeneratingQR {
		session.mu.Unlock()
		return
	}
	session.IsGeneratingQR = true
	session.mu.Unlock()

	defer func() {
		session.mu.Lock()
		session.IsGeneratingQR = false
		session.mu.Unlock()
	}()

	ctx := context.Background()
	qrChan, err := session.Client.GetQRChannel(ctx)
	if err != nil {
		log.Printf("[%s] Failed to get QR channel: %v", session.ID, err)
		return
	}

	err = session.Client.Connect()
	if err != nil {
		log.Printf("[%s] Failed to connect: %v", session.ID, err)
		return
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			session.mu.Lock()
			session.QRCode = evt.Code
			session.Status = "qr_pending"
			session.mu.Unlock()

			// Update database
			sm.db.Exec("UPDATE wa_sessions SET qr_code = ?, status = ? WHERE id = ?", evt.Code, "qr_pending", session.ID)

			log.Printf("[%s] QR code generated", session.ID)
		} else if evt.Event == "success" {
			session.mu.Lock()
			session.QRCode = ""
			session.Status = "connected"
			now := time.Now()
			session.LastConnected = &now
			session.mu.Unlock()

			// Detect phone number
			if session.Client.Store.ID != nil {
				session.mu.Lock()
				session.PhoneNumber = session.Client.Store.ID.User
				session.mu.Unlock()
			}

			// Update database
			sm.db.Exec("UPDATE wa_sessions SET qr_code = ?, status = ?, phone_number = ?, last_connected_at = ? WHERE id = ?",
				"", "connected", session.PhoneNumber, now, session.ID)

			log.Printf("[%s] Successfully logged in! Phone: %s", session.ID, session.PhoneNumber)
			break
		} else if evt.Event == "timeout" {
			log.Printf("[%s] QR code timed out, regenerating...", session.ID)
			time.Sleep(3 * time.Second)
			go sm.generateQRCode(session)
			break
		} else if evt.Event == "err-client-outdated" {
			log.Printf("[%s] WhatsApp client/session data outdated, resetting local device store...", session.ID)
			go sm.resetOutdatedSession(session)
			break
		} else {
			log.Printf("[%s] QR channel event: %s", session.ID, evt.Event)
		}
	}
}

// resetOutdatedSession clears local linked-device data and starts a fresh QR pairing.
func (sm *SessionManager) resetOutdatedSession(session *WhatsAppSession) {
	session.Client.Disconnect()

	if err := session.Client.Store.Delete(context.Background()); err != nil {
		log.Printf("[%s] Failed to reset local device store: %v", session.ID, err)
		return
	}

	session.mu.Lock()
	session.Status = "qr_pending"
	session.PhoneNumber = ""
	session.QRCode = ""
	session.mu.Unlock()

	sm.db.Exec("UPDATE wa_sessions SET status = ?, phone_number = ?, qr_code = ?, last_connected_at = ? WHERE id = ?",
		"qr_pending", "", "", nil, session.ID)

	log.Printf("[%s] Local session reset complete, generating new QR...", session.ID)
	go sm.generateQRCode(session)
}

// createEventHandler creates event handler for session
func (sm *SessionManager) createEventHandler(session *WhatsAppSession) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.ClientOutdated:
			log.Printf("[%s] Received ClientOutdated event (405): %v", session.ID, v)
			go sm.resetOutdatedSession(session)
			return
		case *events.ConnectFailure:
			if v.Reason == events.ConnectFailureClientOutdated {
				log.Printf("[%s] Connect failure due to outdated client/session: %s", session.ID, v.Reason.String())
				go sm.resetOutdatedSession(session)
				return
			}
		}

		eventType := fmt.Sprintf("%T", evt)

		switch {
		case strings.Contains(eventType, "Connected"):
			session.mu.Lock()
			session.Status = "connected"
			now := time.Now()
			session.LastConnected = &now
			session.QRCode = ""

			if session.Client.Store.ID != nil {
				session.PhoneNumber = session.Client.Store.ID.User
			}
			session.mu.Unlock()

			sm.db.Exec("UPDATE wa_sessions SET status = ?, phone_number = ?, last_connected_at = ?, qr_code = ? WHERE id = ?",
				"connected", session.PhoneNumber, now, "", session.ID)

			log.Printf("[%s] Connected! Phone: %s", session.ID, session.PhoneNumber)

		case strings.Contains(eventType, "Disconnected"):
			session.mu.Lock()
			session.Status = "disconnected"
			session.mu.Unlock()

			sm.db.Exec("UPDATE wa_sessions SET status = ? WHERE id = ?", "disconnected", session.ID)
			log.Printf("[%s] Disconnected", session.ID)

			// Try to reconnect
			go func() {
				time.Sleep(2 * time.Second)
				if session.Client.Store.ID != nil {
					log.Printf("[%s] Attempting to reconnect...", session.ID)
					if err := session.Client.Connect(); err != nil {
						log.Printf("[%s] Reconnection failed: %v", session.ID, err)
						sm.generateQRCode(session)
					}
				} else {
					sm.generateQRCode(session)
				}
			}()

		case strings.Contains(eventType, "LoggedOut"):
			log.Printf("[%s] Logged out", session.ID)
			go sm.generateQRCode(session)
		}
	}
}

// createReceiptHandler creates receipt handler for session
func (sm *SessionManager) createReceiptHandler(session *WhatsAppSession) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Receipt:
			phone := v.Chat.User
			var status string

			switch v.Type {
			case events.ReceiptTypeDelivered:
				status = "delivered"
			case events.ReceiptTypeRead:
				status = "read"
			case events.ReceiptTypeSender:
				status = "sent"
			default:
				return
			}

			log.Printf("[%s] Receipt: %s - %s", session.ID, phone, status)
			sm.sendWebhook(phone, status)
		}
	}
}

// sendWebhook sends webhook notification
func (sm *SessionManager) sendWebhook(phone, status string) {
	if sm.webhookURL == "" {
		return
	}

	payload := map[string]string{
		"phone":     phone,
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(payload)

	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("POST", sm.webhookURL, strings.NewReader(string(jsonData)))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}()
}

// GetSession gets session by ID
func (sm *SessionManager) GetSession(id string) (*WhatsAppSession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	session, ok := sm.sessions[id]
	return session, ok
}

// ListSessions lists all sessions
func (sm *SessionManager) ListSessions() []*WhatsAppSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*WhatsAppSession, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// UpdateSessionName updates session name
func (sm *SessionManager) UpdateSessionName(id, name string) error {
	session, ok := sm.GetSession(id)
	if !ok {
		return fmt.Errorf("session not found")
	}

	session.mu.Lock()
	session.Name = name
	session.mu.Unlock()

	_, err := sm.db.Exec("UPDATE wa_sessions SET name = ? WHERE id = ?", name, id)
	return err
}

// LogoutSession logs out session
func (sm *SessionManager) LogoutSession(id string) error {
	session, ok := sm.GetSession(id)
	if !ok {
		return fmt.Errorf("session not found")
	}

	if session.Client.IsConnected() {
		session.Client.Disconnect()
	}

	session.mu.Lock()
	session.Status = "disconnected"
	session.mu.Unlock()

	sm.db.Exec("UPDATE wa_sessions SET status = ? WHERE id = ?", "disconnected", id)
	log.Printf("[%s] Logged out", id)

	return nil
}

// DeleteSession deletes session
func (sm *SessionManager) DeleteSession(id string) error {
	session, ok := sm.GetSession(id)
	if !ok {
		return fmt.Errorf("session not found")
	}

	// Logout first
	if session.Client.IsConnected() {
		session.Client.Disconnect()
	}

	// Close container
	if session.Container != nil {
		session.Container.Close()
	}

	// Delete from database
	sm.db.Exec("DELETE FROM wa_sessions WHERE id = ?", id)

	// Remove from memory
	sm.mu.Lock()
	delete(sm.sessions, id)
	sm.mu.Unlock()

	// Delete session directory
	sessionDir := filepath.Join("session", id)
	os.RemoveAll(sessionDir)

	log.Printf("[%s] Session deleted", id)
	return nil
}

// GetConnectedSessions gets connected sessions from list
func (sm *SessionManager) GetConnectedSessions(sessionIDs []string) []*WhatsAppSession {
	if len(sessionIDs) == 0 {
		return nil
	}

	connected := make([]*WhatsAppSession, 0)
	for _, id := range sessionIDs {
		session, ok := sm.GetSession(id)
		if ok && session.Client.IsConnected() && session.Client.Store.ID != nil {
			connected = append(connected, session)
		}
	}
	return connected
}
