package controllers

import (
	"ppw-uas2526-11-rest-api-banking-mini/configs"
	"ppw-uas2526-11-rest-api-banking-mini/models"
	"ppw-uas2526-11-rest-api-banking-mini/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE (POST) - Membuat akun baru untuk user yang sudah ada
func CreateAccount(c *gin.Context) {
    userID := c.GetUint("user_id") // Dari Token Middleware

    // Generate Nomor Rekening Unik
    newAccNumber := utils.GenerateAccountNumber(configs.DB)

    newAccount := models.Account{
        UserID:        userID,
        AccountNumber: newAccNumber,
        Balance:       0,
    }

    if err := configs.DB.Create(&newAccount).Error; err != nil {
        c.JSON(500, gin.H{"error": "Gagal membuat akun"})
        return
    }

    c.JSON(201, gin.H{
        "message": "Rekening berhasil dibuat",
        "data":    newAccount,
    })
}

// READ (GET) - Mengambil semua akun
func GetAccounts(c *gin.Context) {
	var accounts []models.Account
	configs.DB.Find(&accounts)

	// HTTP 200: OK
	c.JSON(http.StatusOK, gin.H{"data": accounts})
}

// READ (GET by ID) - Mengambil satu akun spesifik
func GetAccountByID(c *gin.Context) {
	var account models.Account
	id := c.Param("id")

	if err := configs.DB.First(&account, id).Error; err != nil {
		// HTTP 404: Not Found jika ID tidak ada
		c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account})
}

// UPDATE (PUT) - Mengupdate data akun
func UpdateAccount(c *gin.Context) {
	var account models.Account
	id := c.Param("id")

	// Cek apakah data ada
	if err := configs.DB.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
		return
	}

	// Validasi Input Update
	var input models.Account
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data
	configs.DB.Model(&account).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": account})
}

// DELETE (DELETE) - Menghapus akun
func DeleteAccount(c *gin.Context) {
	var account models.Account
	id := c.Param("id")

	if err := configs.DB.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
		return
	}

	configs.DB.Delete(&account)

	c.JSON(http.StatusOK, gin.H{"message": "Akun berhasil di hapus"})
}