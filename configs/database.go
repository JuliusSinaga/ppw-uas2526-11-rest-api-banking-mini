package configs

import (
	"fmt"
	"log"
	"os"
	"ppw-uas2526-11-rest-api-banking-mini/models" 

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Setup TimeZone Asia/Jakarta (WIB)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port,
	)

	fmt.Println("⏳ Mencoba koneksi ke database...") // Log tambahan untuk debug

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke PostgreSQL: %v", err)
	}

	// Auto Migrate kedua model
	err = database.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{})
	if err != nil {
		panic("Gagal migrasi database!")
	}

	fmt.Println("✅ Database connected successfully!")
	DB = database
}