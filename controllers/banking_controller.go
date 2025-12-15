package controllers

import (
	"ppw-uas2526-11-rest-api-banking-mini/configs"
	"ppw-uas2526-11-rest-api-banking-mini/models"
	"ppw-uas2526-11-rest-api-banking-mini/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// REGISTER (POST /api/register)
func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash Password sebelum simpan
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	// Buat User & Akun otomatis (Saldo 0)
	err := configs.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&input).Error; err != nil {
			return err
		}

		newAccNumber := utils.GenerateAccountNumber(tx)

        // Buat akun bank default
        account := models.Account{
            UserID:        input.ID,
            AccountNumber: newAccNumber, // Pakai nomor hasil generate
            Balance:       0,
        }
        
        if err := tx.Create(&account).Error; err != nil {
            return err
        }
        return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal register: " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "User berhasil didaftarkan"})
}

// LOGIN (POST /api/login)
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Input tidak valid"})
		return
	}

	var user models.User
	if err := configs.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Email atau password salah"})
		return
	}

	// Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Email atau password salah"})
		return
	}

	// Generate JWT
	token, _ := utils.GenerateToken(user.ID)
	c.JSON(200, gin.H{"token": token})
}

// Helper: Ambil Account ID dari Token (Context)
func getAccountIDFromContext(c *gin.Context) (models.Account, error) {
	userID := c.GetUint("user_id") // Diset oleh Middleware nanti
	var account models.Account
	if err := configs.DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return account, err
	}
	return account, nil
}

// TOPUP (POST /api/topup)
func Topup(c *gin.Context) {
	var input struct {
		Amount float64 `json:"amount" binding:"required,min=10000"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	account, err := getAccountIDFromContext(c)
	if err != nil {
		c.JSON(404, gin.H{"error": "Akun tidak ditemukan"})
		return
	}

	// Tambah Saldo & Catat Mutasi
	configs.DB.Transaction(func(tx *gorm.DB) error {
		account.Balance += input.Amount
		tx.Save(&account)

		tx.Create(&models.Transaction{
			AccountID:       account.ID,
			TransactionType: "TOPUP",
			Amount:          input.Amount,
			Description:     "Topup saldo via API",
		})
		return nil
	})

	c.JSON(200, gin.H{"message": "Topup berhasil", "balance": account.Balance})
}

// TRANSFER (POST /api/transfer)
func Transfer(c *gin.Context) {
	var input struct {
		ToAccountID uint    `json:"to_account_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,min=1000"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sender, err := getAccountIDFromContext(c) // Pengirim diambil dari Token Login
	if err != nil { return } // Error sudah dihandle di helper, sesuaikan logic error handling jika perlu

	if sender.Balance < input.Amount {
		c.JSON(400, gin.H{"error": "Saldo tidak mencukupi"})
		return
	}

	err = configs.DB.Transaction(func(tx *gorm.DB) error {
		var receiver models.Account
		if err := tx.First(&receiver, input.ToAccountID).Error; err != nil {
			return err
		}

		// Update Saldo
		sender.Balance -= input.Amount
		receiver.Balance += input.Amount
		tx.Save(&sender)
		tx.Save(&receiver)

		// Catat Mutasi Pengirim (Uang Keluar)
		tx.Create(&models.Transaction{
			AccountID:       sender.ID,
			TargetAccountID: &receiver.ID,
			TransactionType: "TRANSFER_OUT",
			Amount:          input.Amount,
			Description:     "Transfer keluar",
		})

		// Catat Mutasi Penerima (Uang Masuk)
		tx.Create(&models.Transaction{
			AccountID:       receiver.ID,
			TargetAccountID: &sender.ID,
			TransactionType: "TRANSFER_IN",
			Amount:          input.Amount,
			Description:     "Transfer masuk",
		})

		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Transfer gagal"})
		return
	}

	c.JSON(200, gin.H{"message": "Transfer berhasil"})
}

// CEK SALDO (GET /api/balance)
func GetBalance(c *gin.Context) {
	account, err := getAccountIDFromContext(c)
	if err != nil {
		c.JSON(404, gin.H{"error": "Akun tidak ditemukan"})
		return
	}
	c.JSON(200, gin.H{"balance": account.Balance})
}

// MUTASI (GET /api/mutations)
func GetMutations(c *gin.Context) {
	account, err := getAccountIDFromContext(c)
	if err != nil {
		c.JSON(404, gin.H{"error": "Akun tidak ditemukan"})
		return
	}

	var transactions []models.Transaction
	configs.DB.Where("account_id = ?", account.ID).Order("created_at desc").Find(&transactions)

	c.JSON(200, gin.H{"data": transactions})
}

// WITHDRAW (POST /api/withdraw) - tambahan fitur
func Withdraw(c *gin.Context) {
	var input struct {
		Amount float64 `json:"amount" binding:"required,min=10000"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	account, err := getAccountIDFromContext(c)
	if err != nil {
		c.JSON(404, gin.H{"error": "Akun tidak ditemukan"})
		return
	}
	if account.Balance < input.Amount {
		c.JSON(400, gin.H{"error": "Saldo tidak mencukupi"})
		return
	}
	// Kurangi Saldo & Catat Mutasi
	configs.DB.Transaction(func(tx *gorm.DB) error {
		account.Balance -= input.Amount
		tx.Save(&account)
		tx.Create(&models.Transaction{
			AccountID:       account.ID,
			TransactionType: "WITHDRAW",
			Amount:          input.Amount,
			Description:     "Withdraw saldo via API",
		})
		return nil
	})

	c.JSON(200, gin.H{"message": "Withdraw berhasil", "balance": account.Balance})
}