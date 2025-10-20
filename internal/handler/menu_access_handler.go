package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MenuAccessHandler struct {
	svc      service.MenuAccessService
	validate *validator.Validate
}

func NewMenuAccessHandler(svc service.MenuAccessService) *MenuAccessHandler {
	return &MenuAccessHandler{svc: svc, validate: validator.New()}
}

func (h *MenuAccessHandler) GetTreeByRole(c *fiber.Ctx) error {
	roleID, err := parseUint(c.Params("roleId"))
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role id"}) }

	tree, err := h.svc.GetTreeByRole(roleID)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"menus": tree})
}

func (h *MenuAccessHandler) UpdateBulk(c *fiber.Ctx) error {
	roleID, err := parseUint(c.Params("roleId"))
	if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role id"}) }

	var req dto.UpdateMenuAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.svc.UpdateBulk(roleID, req.Items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "akses berhasil disimpan"})
}

func parseUint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 0)
	return uint(u64), err
}
