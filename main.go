package main

import (
	"ppw-uas2526-11-rest-api-banking-mini/configs"
	"ppw-uas2526-11-rest-api-banking-mini/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
	"os"
	"fmt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env, menggunakan environment system jika ada.")
	}

	// 1. Koneksi Database
	configs.ConnectDatabase()

	// Inisialisasi Gin Router
	r := gin.Default()
	// Konfigurasi ini mengizinkan frontend (React) mengakses backend
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Saat development aman, saat production sebaiknya spesifik domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 2. Setup Routes
	routes.SetupRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Selamat datang di API Banking Mini! 💸",
		})
	})

	// 3. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("✅ Server berjalan di http://localhost:" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}