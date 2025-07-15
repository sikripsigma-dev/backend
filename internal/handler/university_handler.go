package handler

import (
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UniversityHandler struct {
	service service.UniversityService
}

func NewUniversityHandler(service service.UniversityService) *UniversityHandler {
	return &UniversityHandler{service}
}

// func (h *UniversityHandler) RegisterRoutes(router fiber.Router) {
// 	uni := router.Group("/universities")
// 	uni.Get("/", h.GetAll)
// 	uni.Get("/:id", h.GetByID)
// 	uni.Post("/", h.Create)
// 	uni.Put("/:id", h.Update)
// 	uni.Delete("/:id", h.Delete)
// }

func (h *UniversityHandler) GetAll(c *fiber.Ctx) error {
	unis, err := h.service.GetAll(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(unis)
}

func (h *UniversityHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	uni, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Universitas tidak ditemukan")
	}
	return c.JSON(uni)
}

type universityRequest struct {
	Name string `json:"name"`
}

func (h *UniversityHandler) Create(c *fiber.Ctx) error {
	var req universityRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Nama universitas wajib diisi")
	}
	uni, err := h.service.Create(c.Context(), req.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(uni)
}

func (h *UniversityHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req universityRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Nama universitas wajib diisi")
	}
	err := h.service.Update(c.Context(), id, req.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *UniversityHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
