package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole defines the available roles in the system
type UserRole string

const (
	RolePemohon UserRole = "pemohon"
	RolePenilai UserRole = "penilai"
)

// User represents a registered user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Never expose password in JSON
	Role      UserRole       `gorm:"size:20;not null;index" json:"role"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Company   string         `gorm:"size:150" json:"company"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RegisterDTO is the request body for user registration
type RegisterDTO struct {
	Name     string   `json:"name" binding:"required,min=3,max=100"`
	Email    string   `json:"email" binding:"required,email,max=100"`
	Password string   `json:"password" binding:"required,min=8,max=128"`
	Role     UserRole `json:"role" binding:"required,oneof=pemohon penilai"`
	Phone    string   `json:"phone" binding:"max=20"`
	Company  string   `json:"company" binding:"max=150"`
}

// LoginDTO is the request body for user login
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is the sanitized user response (no password)
type UserResponse struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
	Phone     string   `json:"phone"`
	Company   string   `json:"company"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User model to a sanitized UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Phone:     u.Phone,
		Company:   u.Company,
		CreatedAt: u.CreatedAt,
	}
}
