package controllers

import (
	"net/http"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type DocumentTypeController struct{}

func NewDocumentTypeController() *DocumentTypeController {
	return &DocumentTypeController{}
}

// GetDocumentTypes handles GET /api/v1/document-types
// Returns list of document types. Public/Pemohon gets only active ones, Admin gets all.
func (dtc *DocumentTypeController) GetDocumentTypes(c *gin.Context) {
	var documentTypes []models.DocumentType
	query := config.DB.Model(&models.DocumentType{})

	// Filter by search query if provided
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(name ILIKE ? OR code ILIKE ? OR requirement ILIKE ? OR description ILIKE ?)", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Filter active if requested
	if c.Query("active_only") == "true" {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("name asc").Find(&documentTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   documentTypes,
		"total":  len(documentTypes),
	})
}

// GetDocumentTypeByID handles GET /api/v1/document-types/:id
func (dtc *DocumentTypeController) GetDocumentTypeByID(c *gin.Context) {
	id := c.Param("id")
	var docType models.DocumentType

	if err := config.DB.First(&docType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document type not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   docType,
	})
}

// CreateDocumentType handles POST /api/v1/document-types (Admin only)
func (dtc *DocumentTypeController) CreateDocumentType(c *gin.Context) {
	var dto models.CreateDocumentTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid document type data",
			"details": err.Error(),
		})
		return
	}

	// Check if code already exists
	var existing models.DocumentType
	if err := config.DB.Where("code = ?", dto.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Document type code already exists",
		})
		return
	}

	isActive := true
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}

	docType := models.DocumentType{
		Code:        dto.Code,
		Name:        dto.Name,
		Requirement: dto.Requirement,
		Description: dto.Description,
		IsActive:    isActive,
	}

	if err := config.DB.Create(&docType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document type"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Document type master created successfully",
		"data":    docType,
	})
}

// UpdateDocumentType handles PUT /api/v1/document-types/:id (Admin only)
func (dtc *DocumentTypeController) UpdateDocumentType(c *gin.Context) {
	id := c.Param("id")
	var docType models.DocumentType

	if err := config.DB.First(&docType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document type not found"})
		return
	}

	var dto models.CreateDocumentTypeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type data", "details": err.Error()})
		return
	}

	docType.Code = dto.Code
	docType.Name = dto.Name
	docType.Requirement = dto.Requirement
	docType.Description = dto.Description
	if dto.IsActive != nil {
		docType.IsActive = *dto.IsActive
	}

	if err := config.DB.Save(&docType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Document type updated successfully",
		"data":    docType,
	})
}

// DeleteDocumentType handles DELETE /api/v1/document-types/:id (Admin only)
func (dtc *DocumentTypeController) DeleteDocumentType(c *gin.Context) {
	id := c.Param("id")
	var docType models.DocumentType

	if err := config.DB.First(&docType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document type not found"})
		return
	}

	if err := config.DB.Delete(&docType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Document type deleted successfully",
	})
}
