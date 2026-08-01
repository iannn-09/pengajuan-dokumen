package models

import (
	"time"

	"gorm.io/gorm"
)

// ProjectStatus defines the workflow states for a project (permohonan dokumen)
type ProjectStatus string

const (
	StatusDraft       ProjectStatus = "DRAFT"
	StatusSubmitted   ProjectStatus = "SUBMITTED"
	StatusUnderReview ProjectStatus = "UNDER_REVIEW"
	StatusRevision    ProjectStatus = "REVISION"
	StatusApproved    ProjectStatus = "APPROVED"
	StatusRejected    ProjectStatus = "REJECTED"
)

// Project represents a document request submission (permohonan dokumen kelayakan)
type Project struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProjectNumber string         `gorm:"size:50;uniqueIndex;not null" json:"project_number"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	CompanyName   string         `gorm:"size:200" json:"company_name"`
	Unit          string         `gorm:"size:200" json:"unit"` // Unit Kerja / Divisi / Pabrik

	// Relation to Master DocumentType
	DocumentTypeID *uint         `gorm:"index" json:"document_type_id,omitempty"`
	DocumentType   *DocumentType `gorm:"foreignKey:DocumentTypeID" json:"document_type,omitempty"`

	// Foreign key to User (pemohon)
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Workflow status
	Status      ProjectStatus `gorm:"size:20;default:'DRAFT';index" json:"status"`
	SubmittedAt *time.Time    `json:"submitted_at,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Documents       []Document       `gorm:"foreignKey:ProjectID" json:"documents,omitempty"`
	ReviewHistories []ReviewHistory  `gorm:"foreignKey:ProjectID" json:"review_histories,omitempty"`
}

// ────────────────────── DTOs ──────────────────────

// CreateProjectDTO is the request body for creating a new project
type CreateProjectDTO struct {
	Title          string `json:"title" binding:"required,min=5,max=255"`
	Description    string `json:"description" binding:"max=5000"`
	CompanyName    string `json:"company_name" binding:"max=200"`
	Unit           string `json:"unit" binding:"max=200"`
	DocumentTypeID *uint  `json:"document_type_id"`
}

// UpdateProjectDTO is the request body for updating a draft/revision project
type UpdateProjectDTO struct {
	Title          string `json:"title" binding:"required,min=5,max=255"`
	Description    string `json:"description" binding:"max=5000"`
	CompanyName    string `json:"company_name" binding:"max=200"`
	Unit           string `json:"unit" binding:"max=200"`
	DocumentTypeID *uint  `json:"document_type_id"`
}

// ProjectListResponse is the paginated response for project listing
type ProjectListResponse struct {
	Data       []Project `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}
