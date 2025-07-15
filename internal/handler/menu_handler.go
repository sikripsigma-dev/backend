package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/models"
	"Skripsigma-BE/internal/service"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler(menuService *service.MenuService) *MenuHandler{
	return &MenuHandler{menuService}
}

func (h *MenuHandler) CreateMenu(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	menu, err := h.menuService.CreateMenu(req)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Menu created successfully",
		"menu": menu,
	})
}

func (h *MenuHandler) UpdateMenu(c *fiber.Ctx) error {

	menu := c.Locals("menu").(*models.Menu)

	var req dto.UpdateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// updatedMenu, err := h.menuService.UpdateMenu(fmt.Sprintf("%d", menu.ID), req)
	updatedMenu, err := h.menuService.UpdateMenu(fmt.Sprintf("%d", menu.ID), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu updated successfully",
		"menu":    updatedMenu,
	})



}

// func (h *MenuHandler) GetAllMenu(c *fiber.Ctx) error {
// 	menu, err := h.menuService.GetAllMenu()

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error})
// 	}

// 	tree := buildMenuTree(menu, nil)

// 	return c.JSON(fiber.Map{"menus": tree})
// }

func (h *MenuHandler) GetAllMenu(c *fiber.Ctx) error {
	rawUser := c.Locals("user")
	if rawUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user not found in context",
		})
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - invalid user type",
		})
	}

	menu, err := h.menuService.GetMenuByRole(user.RoleId, user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tree := buildMenuTree(menu, nil)
	return c.JSON(fiber.Map{"menus": tree})
}


func buildMenuTree(menus []models.Menu, parentID *uint) []fiber.Map {
	var tree []fiber.Map

	for _, m := range menus {
		if (m.ParentID == nil && parentID == nil) || (m.ParentID != nil && parentID != nil && *m.ParentID == *parentID) {
			children := buildMenuTree(menus, &m.ID)

			item := fiber.Map{
				"id":        m.ID,
				"name":      m.Nama,
				"url":       m.URL,
				"icon":		 m.Icon,
				"is_active": m.IsActive,
			}

			if len(children) > 0 {
				item["children"] = children
			}

			tree = append(tree, item)
		}
	}

	return tree
}