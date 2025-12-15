package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"` // Foreign Key
	AccountNumber string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"account_number" binding:"required,numeric,min=10,max=16"` // Validasi: Harus angka, 10-16 digit
	Balance       float64        `gorm:"type:decimal(15,2);default:0" json:"balance" binding:"min=0"` // Validasi: Saldo tidak boleh negatif
	User          *User          `gorm:"foreignKey:UserID" json:"-"` // Asosiasi ke User (disembunyikan di JSON untuk cegah loop)
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}