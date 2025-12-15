package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" binding:"required,min=3"`     // Validasi: Wajib isi, minimal 3 karakter
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email"` // Validasi: Format email valid & unik
	Password  string         `gorm:"type:varchar(255);not null" json:"password,omitempty" binding:"required,min=6"` // Validasi: Password minimal 6 karakter
	Accounts  []Account      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accounts,omitempty"` // Relasi One-to-Many
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft Delete (tidak ditampilkan di JSON)
}