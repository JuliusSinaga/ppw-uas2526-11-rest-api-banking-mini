package utils

import (
	"math/rand"
	"ppw-uas2526-11-rest-api-banking-mini/models"
	"time"

	"gorm.io/gorm"
)

func init() {
	// Inisialisasi seed agar angka acak selalu berbeda tiap run
	rand.Seed(time.Now().UnixNano())
}

// GenerateAccountNumber membuat nomor rekening acak 12 digit
// dan memastikan tidak ada duplikat di database.
func GenerateAccountNumber(db *gorm.DB) string {
	const charset = "0123456789"
	length := 12 // Standar panjang nomor rekening (bisa diubah)

	for {
		// Generate String Acak
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		candidate := string(b)

		// Cek ke Database (Collision Check)
		var count int64
		// Cek apakah nomor ini sudah dipakai di tabel accounts
		err := db.Model(&models.Account{}).Where("account_number = ?", candidate).Count(&count).Error
		
		// Jika error db (jarang terjadi), atau count == 0 (belum dipakai), kembalikan nomor ini
		if err == nil && count == 0 {
			return candidate
		}
		
		// Jika count > 0, loop akan berulang untuk generate nomor baru
	}
}