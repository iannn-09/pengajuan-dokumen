package models

import (
	"time"

	"gorm.io/gorm"
)

// DocumentType represents the master data for document types (Jenis Dokumen)
type DocumentType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"` // E.g. AMDAL, UKL-UPL, SLF
	Name        string         `gorm:"size:255;not null" json:"name"`            // E.g. Analisis Mengenai Dampak Lingkungan
	Requirement string         `gorm:"type:text" json:"requirement"`             // Rincian persyaratan berkas wajib
	Description string         `gorm:"type:text" json:"description"`             // Deskripsi umum & target pengerjaan
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateDocumentTypeDTO is the request payload for creating/updating a DocumentType
type CreateDocumentTypeDTO struct {
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Requirement string `json:"requirement" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}
