package controllers

import (
	"Skripsigma-BE/database"
	"Skripsigma-BE/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Cek apakah email sudah digunakan
	var existingUser models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	}

	// Hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Membuat user baru dengan UUID
	user := models.User{
		Nim:     data["nim"],
		Name:     data["name"],
		Phone:     data["phone"],
		Email:    data["email"],
		Password: string(hashedPassword),
	}

	// Simpan user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	// Return data user tanpa password
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}


func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User

	// Cari user berdasarkan email
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Cek password yang di-hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Buat token dengan RegisteredClaims
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam
	claims := jwt.RegisteredClaims{
		Issuer:    "Skripsigma-API", // Identitas aplikasi yang menerbitkan token
		Subject:   user.Id,          // ID user sebagai subjek token
		ExpiresAt: jwt.NewNumericDate(expirationTime), // Waktu kedaluwarsa token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Tentukan apakah perlu Secure=true (hanya jika HTTPS)
	secure := c.Protocol() == "https"

	// Simpan token di cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    t,
		Expires:  expirationTime,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax", // Sesuai best practice
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   t,
	})
}


func User(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authentication token"})
	}

	// Verifikasi token JWT
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Ambil ID user dari klaim token
	userId := claims.Subject

	// Cari user berdasarkan ID
	var user models.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"message": "User data retrieved successfully",
		"user": fiber.Map{
			"id":    user.Id,
			"nim":	 user.Nim,
			"name":  user.Name,
			"phone": user.Phone,
			"email": user.Email,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired 1 jam yang lalu
		HTTPOnly: true,
		Secure:   c.Protocol() == "https", // Gunakan Secure jika menggunakan HTTPS
		SameSite: "Lax",
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}