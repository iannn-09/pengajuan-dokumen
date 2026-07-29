package controllers

import (
	"math"
	"net/http"
	"strconv"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// GetUsers handles GET /api/v1/users
// Lists all users with pagination (penilai only)
func (uc *UserController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	var total int64
	query := config.DB.Model(&models.User{})

	// Filter by role
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}

	// Search by name or email
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(name ILIKE ? OR email ILIKE ? OR company ILIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	var users []models.User
	offset := (page - 1) * perPage
	if err := query.Order("created_at desc").Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Convert to response DTOs (no password exposed)
	var userResponses []models.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, u.ToResponse())
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   userResponses,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"per_page":    perPage,
			"total_pages": totalPages,
		},
	})
}

// GetUserByID handles GET /api/v1/users/:id
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user.ToResponse(),
	})
}

// UpdateProfile handles PUT /api/v1/auth/profile
// User can update their own profile
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	type UpdateProfileDTO struct {
		Name    string `json:"name" binding:"required,min=3,max=100"`
		Phone   string `json:"phone" binding:"max=20"`
		Company string `json:"company" binding:"max=150"`
	}

	var dto UpdateProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile data", "details": err.Error()})
		return
	}

	user.Name = dto.Name
	user.Phone = dto.Phone
	user.Company = dto.Company

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
		"data":    user.ToResponse(),
	})
}
