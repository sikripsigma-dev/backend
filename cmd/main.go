package main

import (
	"Skripsigma-BE/internal/config"
	"Skripsigma-BE/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Menghubungkan ke database
	config.ConnectDB()

	// Pastikan koneksi database tertutup saat aplikasi berhenti
	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatal("Failed to get DB object", err)
	}
	defer sqlDB.Close()

	// Membuat instance Fiber
	app := fiber.New()

	routes.WebSocketRoute(app)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowCredentials: true,
	}))

	// Setup routes
	routes.Setup(app)
	app.Static("/images", "./public/images")

	// Menjalankan server
	log.Fatal(app.Listen(":3001"))
}