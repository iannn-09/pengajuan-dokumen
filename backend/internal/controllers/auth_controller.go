package controllers

import (
	"net/http"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Register handles POST /api/v1/auth/register
func (ac *AuthController) Register(c *gin.Context) {
	var dto models.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid registration data",
			"details": err.Error(),
		})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already registered",
		})
		return
	}

	// Hash password with bcrypt
	hashedPassword, err := middleware.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process registration",
		})
		return
	}

	user := models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     models.RolePemohon,
		Phone:    dto.Phone,
		Company:  dto.Company,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user account",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Account created but failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registration successful",
		"data": gin.H{
			"user":  user.ToResponse(),
			"token": token,
		},
	})
}

// Login handles POST /api/v1/auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var dto models.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid login data",
			"details": err.Error(),
		})
		return
	}

	// Find user by email — generic error message to prevent user enumeration
	var user models.User
	if err := config.DB.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Verify password against bcrypt hash
	if !middleware.CheckPassword(dto.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"user":  user.ToResponse(),
			"token": token,
		},
	})
}

// Me handles GET /api/v1/auth/me — returns current user profile
func (ac *AuthController) Me(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to identify user",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user.ToResponse(),
	})
}
