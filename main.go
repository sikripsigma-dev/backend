package main

import (
	"Skripsigma-BE/database"
	"Skripsigma-BE/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Menghubungkan ke database
	database.Connect()

	// Pastikan koneksi database tertutup saat aplikasi berhenti
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Failed to get DB object", err)
	}
	defer sqlDB.Close()

	// Membuat instance Fiber
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowCredentials: true,
	}))

	// Setup routes
	routes.Setup(app)

	// Menjalankan server
	log.Fatal(app.Listen(":3001"))
}