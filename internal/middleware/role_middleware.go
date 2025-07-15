package middleware

import (
	"Skripsigma-BE/internal/models"

	"github.com/gofiber/fiber/v2"
)


func StudentOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(models.User)
        if user.RoleId != 3 {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Only students can access this route.",
            })
        }
        return c.Next()
    }
}

func CompanyOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(models.User)
        if user.RoleId != 2 {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Only companies can access this route.",
            })
        }
        return c.Next()
    }
}

func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(models.User)
        if user.RoleId != 1 {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Only admins can access this route.",
            })
        }
        return c.Next()
    }
}