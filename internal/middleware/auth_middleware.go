package middleware

import (
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari cookie
		token := c.Cookies("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication token",
			})
		}

		// Validasi token & ambil user
		user, err := authService.GetUserByToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Simpan user ke context agar bisa dipakai di handler
		c.Locals("user", user)
		return c.Next()
	}
}
