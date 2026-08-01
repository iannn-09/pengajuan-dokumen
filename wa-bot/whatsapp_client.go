package main

import (
	"context"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// WhatsAppSession represents a single WhatsApp connection
type WhatsAppSession struct {
	ID             string
	Name           string
	PhoneNumber    string
	Status         string // "connected", "disconnected", "qr_pending"
	QRCode         string
	Client         *whatsmeow.Client
	Container      *sqlstore.Container
	CreatedAt      time.Time
	LastConnected  *time.Time
	IsGeneratingQR bool
	mu             sync.RWMutex
}

// Helper functions

func normalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)

	if len(phone) > 0 {
		firstChar := phone[0]
		if firstChar == '\'' || firstChar == '"' || firstChar == '`' {
			phone = phone[1:]
			phone = strings.TrimSpace(phone)
		}
	}

	re := regexp.MustCompile(`\D`)
	phone = re.ReplaceAllString(phone, "")

	if strings.HasPrefix(phone, "08") {
		phone = "628" + phone[2:]
	}

	if strings.HasPrefix(phone, "8") && !strings.HasPrefix(phone, "62") {
		phone = "62" + phone
	}

	if strings.HasPrefix(phone, "+62") {
		phone = phone[1:]
	}

	if !strings.HasPrefix(phone, "62") {
		phone = "62" + phone
	}

	return phone
}

func sendMessageWithRetry(client *whatsmeow.Client, targetJID types.JID, message string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = client.SendMessage(context.Background(), targetJID, &waE2E.Message{
			Conversation: proto.String(message),
		})

		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return err
}

func sendMessageWithRetryDocument(client *whatsmeow.Client, targetJID types.JID, message *waE2E.Message, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = client.SendMessage(context.Background(), targetJID, message)
		if err == nil {
			return nil
		}
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return err
}

func sendSuratTugasPDFWithCaption(client *whatsmeow.Client, targetJID types.JID, caption string) error {
	pdfPath := "surat_tugas_ojk.pdf"

	fileBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return sendMessageWithRetry(client, targetJID, caption, 3)
	}

	uploaded, err := client.Upload(context.Background(), fileBytes, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			Mimetype:      proto.String("application/pdf"),
			Title:         proto.String("Surat_Tugas_OJK_Survei_2025.pdf"),
			FileName:      proto.String("Surat_Tugas_OJK_Survei_2025.pdf"),
			FileLength:    proto.Uint64(uint64(len(fileBytes))),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			DirectPath:    proto.String(uploaded.DirectPath),
			Caption:       proto.String(caption),
		},
	}

	return sendMessageWithRetryDocument(client, targetJID, msg, 2)
}
