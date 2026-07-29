package controllers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectController struct{}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

// ─── Pemohon Endpoints ─────────────────────────────────────────────

// CreateProject handles POST /api/v1/projects
// Only pemohon can create projects. New projects start with DRAFT status.
func (pc *ProjectController) CreateProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	var dto models.CreateProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project data", "details": err.Error()})
		return
	}

	// Generate unique project number: PRJ-YYYYMMDD-XXXX
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	projectNumber := fmt.Sprintf("PRJ-%s-%04d", time.Now().Format("20060102"), r.Intn(10000))

	project := models.Project{
		ProjectNumber: projectNumber,
		Title:         dto.Title,
		Description:   dto.Description,
		CompanyName:   dto.CompanyName,
		UserID:        userID,
		Status:        models.StatusDraft,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Reload with user relation
	config.DB.Preload("User").First(&project, project.ID)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Project created successfully",
		"data":    project,
	})
}

// GetMyProjects handles GET /api/v1/projects
// Returns paginated list of projects owned by the authenticated pemohon
func (pc *ProjectController) GetMyProjects(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	var total int64
	query := config.DB.Model(&models.Project{}).Where("user_id = ?", userID)

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Search by title or project_number
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(title ILIKE ? OR project_number ILIKE ? OR company_name ILIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	var projects []models.Project
	offset := (page - 1) * perPage
	if err := query.Preload("User").Order("created_at desc").Offset(offset).Limit(perPage).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projects,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"per_page":    perPage,
			"total_pages": totalPages,
		},
	})
}

// GetProjectByID handles GET /api/v1/projects/:id
// Returns full project detail with documents and review history
func (pc *ProjectController) GetProjectByID(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	id := c.Param("id")
	var project models.Project

	if err := config.DB.
		Preload("User").
		Preload("Documents").
		Preload("ReviewHistories", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Preload("ReviewHistories.Reviewer").
		First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Pemohon can only see their own projects
	role, _ := middleware.GetUserRole(c)
	if role == "pemohon" && project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   project,
	})
}

// UpdateProject handles PUT /api/v1/projects/:id
// Pemohon can edit a project only when it is in DRAFT or REVISION status
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	id := c.Param("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own projects"})
		return
	}

	if project.Status != models.StatusDraft && project.Status != models.StatusRevision {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot edit project with status '%s'. Only DRAFT or REVISION projects can be edited.", project.Status),
		})
		return
	}

	var dto models.UpdateProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project data", "details": err.Error()})
		return
	}

	project.Title = dto.Title
	project.Description = dto.Description
	project.CompanyName = dto.CompanyName

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project updated successfully",
		"data":    project,
	})
}

// SubmitProject handles POST /api/v1/projects/:id/submit
// Transitions project from DRAFT or REVISION → SUBMITTED
func (pc *ProjectController) SubmitProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	id := c.Param("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only submit your own projects"})
		return
	}

	if project.Status != models.StatusDraft && project.Status != models.StatusRevision {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot submit project with status '%s'. Only DRAFT or REVISION projects can be submitted.", project.Status),
		})
		return
	}

	now := time.Now()
	project.Status = models.StatusSubmitted
	project.SubmittedAt = &now

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project submitted for review",
		"data":    project,
	})
}

// DeleteProject handles DELETE /api/v1/projects/:id
// Pemohon can only delete DRAFT projects
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	id := c.Param("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own projects"})
		return
	}

	if project.Status != models.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only DRAFT projects can be deleted",
		})
		return
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project deleted successfully",
	})
}
