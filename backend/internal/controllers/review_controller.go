package controllers

import (
	"math"
	"net/http"
	"strconv"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewController struct{}

func NewReviewController() *ReviewController {
	return &ReviewController{}
}

// GetProjectsForReview handles GET /api/v1/reviews/projects
// Penilai can see all projects with status SUBMITTED or UNDER_REVIEW
func (rc *ReviewController) GetProjectsForReview(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	var total int64
	query := config.DB.Model(&models.Project{})

	// Default: show SUBMITTED and UNDER_REVIEW projects
	statusFilter := c.Query("status")
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	} else {
		query = query.Where("status IN ?", []string{
			string(models.StatusSubmitted),
			string(models.StatusUnderReview),
		})
	}

	// Search
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(title ILIKE ? OR project_number ILIKE ? OR company_name ILIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	var projects []models.Project
	offset := (page - 1) * perPage
	if err := query.
		Preload("User").
		Order("submitted_at desc, created_at desc").
		Offset(offset).Limit(perPage).
		Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects for review"})
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

// GetProjectForReview handles GET /api/v1/reviews/projects/:id
// Penilai can see full detail of any project (with documents and review history)
func (rc *ReviewController) GetProjectForReview(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   project,
	})
}

// AssessProject handles POST /api/v1/reviews/projects/:id/assess
// Penilai can approve, request revision, or reject a project
func (rc *ReviewController) AssessProject(c *gin.Context) {
	reviewerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify reviewer"})
		return
	}

	id := c.Param("id")
	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Project must be in SUBMITTED or UNDER_REVIEW status to be assessed
	if project.Status != models.StatusSubmitted && project.Status != models.StatusUnderReview {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only SUBMITTED or UNDER_REVIEW projects can be assessed",
		})
		return
	}

	var dto models.AssessProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assessment data", "details": err.Error()})
		return
	}

	statusFrom := project.Status
	var statusTo models.ProjectStatus
	var message string

	switch dto.Action {
	case "approve":
		statusTo = models.StatusApproved
		message = "Project approved successfully"
	case "revise":
		statusTo = models.StatusRevision
		message = "Revision requested for project"
		if dto.Notes == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Notes are required when requesting revision"})
			return
		}
	case "reject":
		statusTo = models.StatusRejected
		message = "Project rejected"
		if dto.Notes == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Notes are required when rejecting a project"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Must be 'approve', 'revise', or 'reject'"})
		return
	}

	// Use transaction to update project status and create review history atomically
	tx := config.DB.Begin()

	project.Status = statusTo
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project status"})
		return
	}

	reviewHistory := models.ReviewHistory{
		ProjectID:  project.ID,
		ReviewerID: reviewerID,
		StatusFrom: statusFrom,
		StatusTo:   statusTo,
		Notes:      dto.Notes,
	}

	if err := tx.Create(&reviewHistory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review history"})
		return
	}

	tx.Commit()

	// Reload with reviewer relation
	config.DB.Preload("Reviewer").First(&reviewHistory, reviewHistory.ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data": gin.H{
			"project":        project,
			"review_history": reviewHistory,
		},
	})
}

// GetMyReviewHistory handles GET /api/v1/reviews/history
// Returns all review history entries created by the authenticated penilai
func (rc *ReviewController) GetMyReviewHistory(c *gin.Context) {
	reviewerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify reviewer"})
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
	query := config.DB.Model(&models.ReviewHistory{}).Where("reviewer_id = ?", reviewerID)
	query.Count(&total)

	var histories []models.ReviewHistory
	offset := (page - 1) * perPage
	if err := query.
		Preload("Reviewer").
		Preload("Project").
		Preload("Project.User").
		Order("created_at desc").
		Offset(offset).Limit(perPage).
		Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review history"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   histories,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"per_page":    perPage,
			"total_pages": totalPages,
		},
	})
}

// GetAllReviewHistory handles GET /api/v1/reviews/all-history
// Returns all review history entries across all reviewers (for dashboard)
func (rc *ReviewController) GetAllReviewHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	var total int64
	query := config.DB.Model(&models.ReviewHistory{})

	// Filter by project_id if provided
	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	query.Count(&total)

	var histories []models.ReviewHistory
	offset := (page - 1) * perPage
	if err := query.
		Preload("Reviewer").
		Preload("Project").
		Preload("Project.User").
		Order("created_at desc").
		Offset(offset).Limit(perPage).
		Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review history"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   histories,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"per_page":    perPage,
			"total_pages": totalPages,
		},
	})
}
