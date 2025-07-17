package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"Skripsigma-BE/internal/models"

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

func ConnectDB() {
	LoadEnv()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Database connected successfully!")

	autoMigrate()
}

func autoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.ResearchCase{},
		&models.Tag{},
		&models.ResearchCaseTag{},
		&models.Application{},
		&models.CompanyUser{},
		&models.Role{},
		&models.Menu{},
		&models.MenuAccess{},
		&models.ChatRooms{},
		&models.ChatMessage{},
		&models.Notification{},
		&models.WeeklyReport{},
		&models.University{},
		&models.StudentUser{},
		&models.SupervisorUser{},
		&models.UserCreateLog{},
		&models.Assignment{},
		&models.StudentDocument{},
		&models.CompanyWeeklyReport{},
	)

	if err != nil {
		log.Fatalf("❌ Gagal migrasi: %v", err)
	}
	log.Println("✅ Migrasi database berhasil!")
}
