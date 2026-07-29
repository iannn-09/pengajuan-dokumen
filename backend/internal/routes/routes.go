package routes

import (
	"net/http"
	"os"
	"time"

	"pengajuan-dokumen/backend/internal/controllers"
	"pengajuan-dokumen/backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Set max multipart memory for file uploads (10 MB)
	r.MaxMultipartMemory = 10 << 20

	// CORS Configuration
	allowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
	if allowOrigin == "" {
		allowOrigin = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowOrigin, "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Version 1 Group
	v1 := r.Group("/api/v1")
	{
		// ─── Public Routes ─────────────────────────────────────────
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "Pengajuan Dokumen API",
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		// Auth (public)
		authController := controllers.NewAuthController()
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
		}

		// ─── Protected Routes (JWT Required) ───────────────────────
		protected := v1.Group("")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			// Auth - profile
			userController := controllers.NewUserController()
			protected.GET("/auth/me", authController.Me)
			protected.PUT("/auth/profile", userController.UpdateProfile)

			// ─── Projects (Pemohon) ────────────────────────────────
			projectController := controllers.NewProjectController()
			projectGroup := protected.Group("/projects")
			projectGroup.Use(middleware.RoleMiddleware("pemohon"))
			{
				projectGroup.GET("", projectController.GetMyProjects)
				projectGroup.POST("", projectController.CreateProject)
				projectGroup.GET("/:id", projectController.GetProjectByID)
				projectGroup.PUT("/:id", projectController.UpdateProject)
				projectGroup.POST("/:id/submit", projectController.SubmitProject)
				projectGroup.DELETE("/:id", projectController.DeleteProject)

				// Document upload for projects
				docUploadController := controllers.NewDocumentUploadController()
				projectGroup.POST("/:id/documents", docUploadController.UploadDocument)
				projectGroup.GET("/:id/documents", docUploadController.GetProjectDocuments)
				projectGroup.DELETE("/:id/documents/:docId", docUploadController.DeleteDocument)
			}

			// ─── Reviews (Penilai) ─────────────────────────────────
			reviewController := controllers.NewReviewController()
			reviewGroup := protected.Group("/reviews")
			reviewGroup.Use(middleware.RoleMiddleware("penilai"))
			{
				reviewGroup.GET("/projects", reviewController.GetProjectsForReview)
				reviewGroup.GET("/projects/:id", reviewController.GetProjectForReview)
				reviewGroup.POST("/projects/:id/assess", reviewController.AssessProject)
				reviewGroup.GET("/history", reviewController.GetMyReviewHistory)
				reviewGroup.GET("/all-history", reviewController.GetAllReviewHistory)
			}

			// ─── Document Download (Both roles) ────────────────────
			docUploadController := controllers.NewDocumentUploadController()
			protected.GET("/documents/:id/download", docUploadController.DownloadDocument)

			// ─── Dashboard (Both roles) ────────────────────────────
			dashboardController := controllers.NewDashboardController()
			protected.GET("/dashboard/stats", dashboardController.GetStats)
			protected.GET("/dashboard/chart-data", dashboardController.GetChartData)

			// ─── User Management (Penilai only) ────────────────────
			userGroup := protected.Group("/users")
			userGroup.Use(middleware.RoleMiddleware("penilai"))
			{
				userGroup.GET("", userController.GetUsers)
				userGroup.GET("/:id", userController.GetUserByID)
			}
		}
	}

	return r
}
