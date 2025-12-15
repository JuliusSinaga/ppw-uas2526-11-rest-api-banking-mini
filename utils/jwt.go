package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("kunci_rahasia_saya_123") // Harusnya dari os.Getenv

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// Fungsi untuk validasi token (dipakai nanti di Middleware)
func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	// Konversi user_id dari float64 (default JWT json) ke uint
	userID := uint(claims["user_id"].(float64))
	return userID, nil
}