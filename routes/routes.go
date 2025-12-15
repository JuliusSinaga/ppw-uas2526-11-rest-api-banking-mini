package routes

import (
	"ppw-uas2526-11-rest-api-banking-mini/controllers"
	"ppw-uas2526-11-rest-api-banking-mini/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware Auth untuk memvalidasi Token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		
		// Cek apakah header Authorization ada
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token diperlukan untuk akses ini"})
			c.Abort()
			return
		}

		// Hapus prefix "Bearer " jika ada (standar Authorization header)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// Validasi token menggunakan fungsi dari utils
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Token tidak valid atau kadaluwarsa"})
			c.Abort()
			return
		}

		// Simpan user_id ke dalam context agar bisa dipakai di controller
		c.Set("user_id", userID)
		c.Next()
	}
}

func SetupRoutes(r *gin.Engine) {
	// Group API (sesuai gambar endpoint: /api/...)
	api := r.Group("/api/v1")
	{
		// --- Public Routes ---
		// Register User Baru
		api.POST("/register", controllers.Register)
		// Login & Dapatkan Token
		api.POST("/login", controllers.Login)

		api.GET("/accounts", controllers.GetAccounts)
		api.GET("/accounts/:id", controllers.GetAccountByID)
		// --- Protected Routes ---
		protected := api.Group("/")
		protected.Use(AuthMiddleware())
		{
			// Topup Saldo
			protected.POST("/topup", controllers.Topup)
			// Transfer Antar Akun
			protected.POST("/transfer", controllers.Transfer)
			// Cek Saldo
			protected.GET("/balance", controllers.GetBalance)
			// Cek Riwayat Mutasi
			protected.GET("/mutations", controllers.GetMutations)
			// Withdraw Saldo
			protected.POST("/withdraw", controllers.Withdraw)
			// Dapatkan User beserta Akun-akunnya
			protected.GET("/users/:id/accounts", controllers.GetUserAndAccounts)
			// Buat Akun Baru untuk User yang sudah ada
			protected.POST("/accounts", controllers.CreateAccount)
			// Update Akun
			protected.PUT("/accounts/:id", controllers.UpdateAccount)
			// Hapus Akun
			protected.DELETE("/accounts/:id", controllers.DeleteAccount)
		}
	}
}