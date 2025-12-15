package controllers

import (
	"net/http"
	"ppw-uas2526-11-rest-api-banking-mini/configs"
	"ppw-uas2526-11-rest-api-banking-mini/models"

	"github.com/gin-gonic/gin"
)

// GET: Menampilkan User beserta Akun-akunnya (Eager Loading)
func GetUserAndAccounts(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Preload("Accounts") akan otomatis mengambil data dari tabel accounts yang berelasi
	if err := configs.DB.Preload("Accounts").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}