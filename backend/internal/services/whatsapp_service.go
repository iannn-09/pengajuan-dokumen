package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"pengajuan-dokumen/backend/internal/models"

	_ "github.com/lib/pq"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	WAClient      *whatsmeow.Client
	WAContainer   *sqlstore.Container
	CurrentQRCode string
	QRMutex       sync.RWMutex
	waOnce        sync.Once
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// syncLatestWAVersion fetches the newest WhatsApp Web version dynamically
func syncLatestWAVersion() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	latestVer, err := whatsmeow.GetLatestVersion(ctx, nil)
	if err != nil {
		log.Printf("[WhatsApp Whatsmeow] Warning: Could not fetch latest WA version: %v\n", err)
		return
	}

	currentVer := store.GetWAVersion()
	if currentVer.LessThan(*latestVer) {
		store.SetWAVersion(*latestVer)
		log.Printf("[WhatsApp Whatsmeow] Updated WhatsApp Web version to: %s\n", latestVer.String())
	}
}

// InitWhatsAppService initializes native whatsmeow WhatsApp client using PostgreSQL database
func InitWhatsAppService() {
	waOnce.Do(func() {
		ctx := context.Background()

		// Configure Device Properties for WhatsApp Multi-Device Compatibility
		store.DeviceProps.Os = proto.String("Windows")
		store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()
		store.DeviceProps.RequireFullSync = proto.Bool(false)

		// Sync latest WhatsApp Web version
		syncLatestWAVersion()

		dbLog := waLog.Stdout("Database", "WARN", true)

		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "dokumen")
		sslmode := getEnv("DB_SSLMODE", "disable")

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)

		container, err := sqlstore.New(ctx, "postgres", dsn, dbLog)
		if err != nil {
			log.Printf("[WhatsApp Whatsmeow] Error initializing PostgreSQL store: %v\n", err)
			return
		}
		WAContainer = container

		deviceStore, err := container.GetFirstDevice(ctx)
		if err != nil {
			log.Printf("[WhatsApp Whatsmeow] Error getting first device from PostgreSQL: %v\n", err)
			return
		}

		clientLog := waLog.Stdout("Whatsmeow", "INFO", true)
		client := whatsmeow.NewClient(deviceStore, clientLog)
		WAClient = client

		if client.Store.ID == nil {
			// Device is not logged in, generate QR Code channel
			go startQRListener(client)
		} else {
			// Device is already paired/logged in
			err := client.Connect()
			if err != nil {
				log.Printf("[WhatsApp Whatsmeow] Error connecting existing session: %v\n", err)
			} else {
				log.Println("[WhatsApp Whatsmeow] Successfully reconnected active WhatsApp session from PostgreSQL!")
			}
		}
	})
}

func startQRListener(client *whatsmeow.Client) {
	ctx := context.Background()
	qrChan, err := client.GetQRChannel(ctx)
	if err != nil {
		if !client.IsConnected() {
			log.Printf("[WhatsApp Whatsmeow] Error getting QR channel: %v\n", err)
		}
		return
	}

	err = client.Connect()
	if err != nil {
		log.Printf("[WhatsApp Whatsmeow] Error connecting client for QR pairing: %v\n", err)
		return
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			// Convert QR string to Base64 Data URL Image
			pngBytes, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err == nil {
				base64Img := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
				QRMutex.Lock()
				CurrentQRCode = base64Img
				QRMutex.Unlock()
				log.Println("[WhatsApp Whatsmeow] New WhatsApp QR Code generated for scanning!")
			}
		} else if evt.Event == "success" || evt.Event == "pair-success" {
			QRMutex.Lock()
			CurrentQRCode = ""
			QRMutex.Unlock()
			log.Println("[WhatsApp Whatsmeow] WhatsApp pairing successful! Client connected.")
		} else if evt.Event == "timeout" {
			log.Println("[WhatsApp Whatsmeow] QR code timeout, refreshing...")
		}
	}
}

// GetWhatsAppStatus returns current native whatsmeow connection status
func GetWhatsAppStatus() (connected bool, phone string, qrCode string) {
	if WAClient == nil {
		return false, "", ""
	}

	connected = WAClient.IsConnected() && WAClient.IsLoggedIn()
	if WAClient.Store.ID != nil {
		phone = WAClient.Store.ID.User
	}

	QRMutex.RLock()
	qrCode = CurrentQRCode
	QRMutex.RUnlock()

	return connected, phone, qrCode
}

// DisconnectWhatsApp logs out current WhatsApp device session in PostgreSQL
func DisconnectWhatsApp() error {
	if WAClient == nil {
		return fmt.Errorf("WhatsApp client not initialized")
	}

	ctx := context.Background()
	err := WAClient.Logout(ctx)
	if err != nil {
		WAClient.Disconnect()
	}

	QRMutex.Lock()
	CurrentQRCode = ""
	QRMutex.Unlock()

	// Restart QR Listener for new connection
	go startQRListener(WAClient)
	return nil
}

// NormalizePhoneNumber formats Indonesian phone numbers to '628...' format
func NormalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "+", "")

	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}
	return phone
}

// SendWhatsappNotification sends a native WhatsApp message via whatsmeow asynchronously
func SendWhatsappNotification(targetPhone string, message string) {
	enabled := os.Getenv("WHATSAPP_NOTIFIER_ENABLED")
	if enabled == "false" {
		return
	}

	if WAClient == nil || !WAClient.IsConnected() || !WAClient.IsLoggedIn() {
		log.Println("[WhatsApp Whatsmeow] Skipping send: Client not logged in or disconnected.")
		return
	}

	phoneFormatted := NormalizePhoneNumber(targetPhone)
	if phoneFormatted == "" {
		log.Println("[WhatsApp Whatsmeow] Warning: Target phone number is empty.")
		return
	}

	recipientJID := types.NewJID(phoneFormatted, types.DefaultUserServer)

	// Run in background goroutine to avoid blocking HTTP API response
	go func() {
		_, err := WAClient.SendMessage(context.Background(), recipientJID, &waProto.Message{
			Conversation: proto.String(message),
		})
		if err != nil {
			log.Printf("[WhatsApp Whatsmeow] Error sending message to %s: %v\n", phoneFormatted, err)
			return
		}
		log.Printf("[WhatsApp Whatsmeow] Native WA message sent successfully to %s\n", phoneFormatted)
	}()
}

// SendStatusChangeNotification formats message with Verifier Name and triggers WA notification for Pemohon
func SendStatusChangeNotification(pemohonName, pemohonPhone, projectNumber, projectTitle, companyName, reviewerName string, statusFrom, statusTo models.ProjectStatus, notes string) {
	if pemohonPhone == "" {
		log.Printf("[WhatsApp Service] Skipping WA notification for %s (%s): phone number missing.\n", pemohonName, projectNumber)
		return
	}

	if reviewerName == "" {
		reviewerName = "Tim Verifikator"
	}

	var statusLabel string
	var statusEmoji string

	switch statusTo {
	case models.StatusApproved:
		statusLabel = "DISETUJUI (APPROVED) ✅"
		statusEmoji = "🎉"
	case models.StatusRevision:
		statusLabel = "PERLU REVISI (REVISION) ⚠️"
		statusEmoji = "📝"
	case models.StatusRejected:
		statusLabel = "DITOLAK (REJECTED) ❌"
		statusEmoji = "🔴"
	case models.StatusUnderReview:
		statusLabel = "SEDANG EVALUASI (UNDER REVIEW) 🔎"
		statusEmoji = "⏳"
	default:
		statusLabel = string(statusTo)
		statusEmoji = "📌"
	}

	message := fmt.Sprintf(
		"Halo *%s*,\n\n"+
			"%s *PEMBERITAHUAN STATUS PERMOHONAN DOKUMEN*\n\n"+
			"Status permohonan dokumen Anda telah diperbarui:\n\n"+
			"📌 *No. Project:* %s\n"+
			"📄 *Judul Project:* %s\n"+
			"🏢 *Perusahaan/Instansi:* %s\n"+
			"👤 *Penilai / Verifikator:* %s\n\n"+
			"Status Terbaru: *%s*\n",
		pemohonName,
		statusEmoji,
		projectNumber,
		projectTitle,
		companyName,
		reviewerName,
		statusLabel,
	)

	if notes != "" {
		message += fmt.Sprintf("\n💬 *Catatan / Evaluasi Penilai:*\n\"%s\"\n", notes)
	}

	message += "\nSilakan login ke portal pengajuan dokumen untuk informasi selengkapnya.\n\n_Pesan otomatis dari Sistem Pengajuan Dokumen Kelayakan._"

	SendWhatsappNotification(pemohonPhone, message)
}
