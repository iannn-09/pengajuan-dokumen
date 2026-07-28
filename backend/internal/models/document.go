package models

import (
	"time"
)

// Document represents an uploaded file attachment belonging to a Project
type Document struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProjectID uint   `gorm:"not null;index" json:"project_id"`
	FileName  string `gorm:"size:255;not null" json:"file_name"`   // Original filename from user
	StoredName string `gorm:"size:255;not null" json:"stored_name"` // UUID-based filename on disk
	FilePath  string `gorm:"size:500;not null" json:"file_path"`   // Relative path within uploads dir
	FileSize  int64  `gorm:"not null" json:"file_size"`            // Size in bytes
	FileType  string `gorm:"size:20" json:"file_type"`             // Extension (pdf, jpg, etc.)
	MimeType  string `gorm:"size:100" json:"mime_type"`            // MIME type (application/pdf, etc.)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
}
