package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Allowed file extensions and max size
var (
	allowedExtensions = map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".doc":  true,
		".docx": true,
	}
	maxFileSize int64 = 10 * 1024 * 1024 // 10 MB
)

type DocumentUploadController struct{}

func NewDocumentUploadController() *DocumentUploadController {
	return &DocumentUploadController{}
}

// UploadDocument handles POST /api/v1/projects/:id/documents
// Allows pemohon to upload files to their project (only DRAFT or REVISION status)
func (duc *DocumentUploadController) UploadDocument(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	projectID := c.Param("id")
	var project models.Project
	if err := config.DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Ownership check
	if project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only upload documents to your own projects"})
		return
	}

	// Status check
	if project.Status != models.StatusDraft && project.Status != models.StatusRevision {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot upload documents to project with status '%s'", project.Status),
		})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file field. Use 'file' as the form field name."})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File size exceeds maximum limit of %d MB", maxFileSize/(1024*1024)),
		})
		return
	}

	// Validate file extension (allow-list approach)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File type '%s' is not allowed. Allowed types: PDF, JPG, JPEG, PNG, DOC, DOCX", ext),
		})
		return
	}

	// Generate unique stored filename using UUID to prevent path traversal
	storedName := uuid.New().String() + ext
	uploadDir := filepath.Join("uploads", fmt.Sprintf("project_%d", project.ID))

	// Create project-specific upload directory
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	storedPath := filepath.Join(uploadDir, storedName)

	// Save file to disk
	if err := c.SaveUploadedFile(header, storedPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	// Sanitize original filename: use only the base name to prevent path components
	sanitizedFileName := filepath.Base(header.Filename)

	// Create document record in database
	doc := models.Document{
		ProjectID:  project.ID,
		FileName:   sanitizedFileName,
		StoredName: storedName,
		FilePath:   storedPath,
		FileSize:   header.Size,
		FileType:   strings.TrimPrefix(ext, "."),
		MimeType:   header.Header.Get("Content-Type"),
	}

	if err := config.DB.Create(&doc).Error; err != nil {
		// Clean up file on DB error
		os.Remove(storedPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Document uploaded successfully",
		"data":    doc,
	})
}

// GetProjectDocuments handles GET /api/v1/projects/:id/documents
func (duc *DocumentUploadController) GetProjectDocuments(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := config.DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Access check for pemohon
	role, _ := middleware.GetUserRole(c)
	if role == "pemohon" {
		userID, _ := middleware.GetUserID(c)
		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view documents of your own projects"})
			return
		}
	}

	var documents []models.Document
	if err := config.DB.Where("project_id = ?", projectID).Order("created_at desc").Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   documents,
		"total":  len(documents),
	})
}

// DeleteDocument handles DELETE /api/v1/projects/:id/documents/:docId
func (duc *DocumentUploadController) DeleteDocument(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to identify user"})
		return
	}

	projectID := c.Param("id")
	docID := c.Param("docId")

	var project models.Project
	if err := config.DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete documents from your own projects"})
		return
	}

	if project.Status != models.StatusDraft && project.Status != models.StatusRevision {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete documents from a submitted project"})
		return
	}

	var doc models.Document
	if err := config.DB.Where("id = ? AND project_id = ?", docID, projectID).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Delete file from disk
	os.Remove(doc.FilePath)

	// Delete record from database
	if err := config.DB.Delete(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Document deleted successfully",
	})
}

// DownloadDocument handles GET /api/v1/documents/:id/download
func (duc *DocumentUploadController) DownloadDocument(c *gin.Context) {
	docID := c.Param("id")

	var doc models.Document
	if err := config.DB.First(&doc, docID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Verify the project exists and check access
	var project models.Project
	if err := config.DB.First(&project, doc.ProjectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Associated project not found"})
		return
	}

	role, _ := middleware.GetUserRole(c)
	if role == "pemohon" {
		userID, _ := middleware.GetUserID(c)
		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only download documents from your own projects"})
			return
		}
	}

	// Verify file exists on disk
	resolvedPath, err := filepath.Abs(doc.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve file path"})
		return
	}

	// Security: verify the resolved path is within the uploads directory
	uploadsDir, _ := filepath.Abs("uploads")
	if !strings.HasPrefix(resolvedPath, uploadsDir+string(filepath.Separator)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid file path"})
		return
	}

	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on disk"})
		return
	}

	// Serve file with security headers
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.FileName))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Type", doc.MimeType)

	// Validate docID is a number to prevent path traversal via parameter
	if _, err := strconv.ParseUint(docID, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	c.File(resolvedPath)
}
