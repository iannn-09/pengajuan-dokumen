package controllers

import (
	"net/http"
	"time"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// GetStats handles GET /api/v1/dashboard/stats
// Returns role-specific statistics
func (dc *DashboardController) GetStats(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}
	role, _ := middleware.GetUserRole(c)

	var stats gin.H

	if role == "pemohon" {
		stats = dc.getPemohonStats(userID)
	} else {
		stats = dc.getPenilaiStats(userID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}

// GetChartData handles GET /api/v1/dashboard/chart-data
// Returns monthly breakdown data for charts
func (dc *DashboardController) GetChartData(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)

	// Monthly submissions for the last 12 months
	type MonthlyData struct {
		Month string `json:"month"`
		Count int64  `json:"count"`
	}

	var monthlySubmissions []MonthlyData
	monthlyQuery := config.DB.Model(&models.Project{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count")

	if role == "pemohon" {
		monthlyQuery = monthlyQuery.Where("user_id = ?", userID)
	}

	monthlyQuery.
		Where("created_at >= ?", time.Now().AddDate(-1, 0, 0)).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month asc").
		Find(&monthlySubmissions)

	// Status distribution
	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var statusDistribution []StatusCount
	statusQuery := config.DB.Model(&models.Project{}).
		Select("status, COUNT(*) as count")

	if role == "pemohon" {
		statusQuery = statusQuery.Where("user_id = ?", userID)
	}

	statusQuery.Group("status").Find(&statusDistribution)

	// Monthly approvals (review history)
	var monthlyApprovals []MonthlyData
	approvalQuery := config.DB.Model(&models.ReviewHistory{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("created_at >= ?", time.Now().AddDate(-1, 0, 0))

	if role == "penilai" {
		approvalQuery = approvalQuery.Where("reviewer_id = ?", userID)
	}

	approvalQuery.
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month asc").
		Find(&monthlyApprovals)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"monthly_submissions":  monthlySubmissions,
			"status_distribution":  statusDistribution,
			"monthly_approvals":    monthlyApprovals,
		},
	})
}

func (dc *DashboardController) getPemohonStats(userID uint) gin.H {
	var total, draft, submitted, underReview, revision, approved, rejected int64

	base := config.DB.Model(&models.Project{}).Where("user_id = ?", userID)
	base.Count(&total)
	config.DB.Model(&models.Project{}).Where("user_id = ? AND status = ?", userID, models.StatusDraft).Count(&draft)
	config.DB.Model(&models.Project{}).Where("user_id = ? AND status = ?", userID, models.StatusSubmitted).Count(&submitted)
	config.DB.Model(&models.Project{}).Where("user_id = ? AND status = ?", userID, models.StatusUnderReview).Count(&underReview)
	config.DB.Model(&models.Project{}).Where("user_id = ? AND status = ?", userID, models.StatusRevision).Count(&revision)
	config.DB.Model(&models.Project{}).Where("user_id = ? AND status = ?", userID, models.StatusApproved).Count(&approved)
	config.DB.Model(&models.Project{}).Where("user_id = ? AND status = ?", userID, models.StatusRejected).Count(&rejected)

	return gin.H{
		"total":        total,
		"draft":        draft,
		"submitted":    submitted,
		"under_review": underReview,
		"revision":     revision,
		"approved":     approved,
		"rejected":     rejected,
	}
}

func (dc *DashboardController) getPenilaiStats(reviewerID uint) gin.H {
	var totalProjects, pendingReview, approved, rejected, revision int64

	config.DB.Model(&models.Project{}).Count(&totalProjects)
	config.DB.Model(&models.Project{}).Where("status IN ?", []string{
		string(models.StatusSubmitted),
		string(models.StatusUnderReview),
	}).Count(&pendingReview)
	config.DB.Model(&models.Project{}).Where("status = ?", models.StatusApproved).Count(&approved)
	config.DB.Model(&models.Project{}).Where("status = ?", models.StatusRejected).Count(&rejected)
	config.DB.Model(&models.Project{}).Where("status = ?", models.StatusRevision).Count(&revision)

	// Count reviews done by this penilai
	var myReviews int64
	config.DB.Model(&models.ReviewHistory{}).Where("reviewer_id = ?", reviewerID).Count(&myReviews)

	// Count total users
	var totalPemohon, totalPenilai int64
	config.DB.Model(&models.User{}).Where("role = ?", models.RolePemohon).Count(&totalPemohon)
	config.DB.Model(&models.User{}).Where("role = ?", models.RolePenilai).Count(&totalPenilai)

	return gin.H{
		"total_projects": totalProjects,
		"pending_review": pendingReview,
		"approved":       approved,
		"rejected":       rejected,
		"revision":       revision,
		"my_reviews":     myReviews,
		"total_pemohon":  totalPemohon,
		"total_penilai":  totalPenilai,
	}
}
