package database

import (
	"Skripsigma-BE/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}
}

func Connect() {
	LoadEnv() // Pastikan .env dimuat

	// Ambil konfigurasi database dari .env
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Koneksi ke database dengan opsi logger
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Menampilkan log query SQL
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Cek koneksi database
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	// Ping database untuk memastikan koneksi
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	log.Println("✅ Database connected successfully!")

	// Set konfigurasi koneksi
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Migrasi semua tabel
	err = DB.AutoMigrate(
		&models.User{},             // ss_users
		&models.Company{},          // ss_m_companies
		&models.ResearchCase{},     // ss_t_research_cases
		&models.Tag{},              // ss_m_tags
		&models.ResearchCaseTag{},  // ss_t_research_case_tags (pivot)
		&models.Application{},      // ss_t_applications
	)

	if err != nil {
		log.Fatalf("❌ Gagal melakukan migrasi database: %v", err)
	}

	log.Println("✅ Migrasi database berhasil!")
}
