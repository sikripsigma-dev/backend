package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService}
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req dto.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	role, err := h.roleService.CreateRole(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Role created successfully",
		"role":    role,
	})
}

func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	
	id := c.Params("id")
	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}

	return c.JSON(fiber.Map{"role": role})
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"roles": roles})
}