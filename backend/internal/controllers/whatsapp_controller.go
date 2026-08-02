package controllers

import (
	"net/http"

	"pengajuan-dokumen/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type WhatsAppController struct{}

func NewWhatsAppController() *WhatsAppController {
	return &WhatsAppController{}
}

type TestWAMessageDTO struct {
	TargetPhone string `json:"target_phone" binding:"required"`
	Message     string `json:"message" binding:"required"`
}

// GetStatus handles GET /api/v1/whatsapp/status
func (wc *WhatsAppController) GetStatus(c *gin.Context) {
	connected, phone, qrCode := services.GetWhatsAppStatus()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"enabled":     true,
			"connected":   connected,
			"phone":       phone,
			"sender_name": "Administrator WhatsApp Gateway (Whatsmeow)",
			"qr_code":     qrCode,
		},
	})
}

// DisconnectSession handles POST /api/v1/whatsapp/disconnect
func (wc *WhatsAppController) DisconnectSession(c *gin.Context) {
	err := services.DisconnectWhatsApp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "WhatsApp session logged out successfully. Scan new QR code to reconnect.",
	})
}

// SendTestMessage handles POST /api/v1/whatsapp/test
func (wc *WhatsAppController) SendTestMessage(c *gin.Context) {
	var dto TestWAMessageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target phone and message are required"})
		return
	}

	services.SendWhatsappNotification(dto.TargetPhone, dto.Message)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Test WhatsApp message sent to " + dto.TargetPhone,
	})
}
