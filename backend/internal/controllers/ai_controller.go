package controllers

import (
	"net/http"

	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AIController struct{}

func NewAIController() *AIController {
	return &AIController{}
}

type AIChatDTO struct {
	Message string `json:"message" binding:"required"`
}

// Chat handles POST /api/v1/ai/chat
func (ac *AIController) Chat(c *gin.Context) {
	var dto AIChatDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	role, _ := middleware.GetUserRole(c)

	reply := services.GenerateAIResponse(dto.Message, role)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"reply": reply,
		},
	})
}
