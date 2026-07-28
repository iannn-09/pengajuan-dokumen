package models

import (
	"time"
)

// ReviewHistory records every status change made by a penilai (reviewer)
// This is the audit log / approval history for a project
type ReviewHistory struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	ProjectID  uint          `gorm:"not null;index" json:"project_id"`
	ReviewerID uint          `gorm:"not null;index" json:"reviewer_id"`
	StatusFrom ProjectStatus `gorm:"size:20;not null" json:"status_from"`
	StatusTo   ProjectStatus `gorm:"size:20;not null" json:"status_to"`
	Notes      string        `gorm:"type:text" json:"notes"` // Catatan penilai

	CreatedAt time.Time `json:"created_at"`

	// Relations
	Reviewer User    `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Project  Project `gorm:"foreignKey:ProjectID" json:"-"`
}

// AssessProjectDTO is the request body for penilai to assess a project
type AssessProjectDTO struct {
	Action string `json:"action" binding:"required,oneof=approve revise reject"`
	Notes  string `json:"notes" binding:"max=5000"`
}
