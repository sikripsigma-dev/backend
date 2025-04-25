package handler

import (
	"Skripsigma-BE/internal/dto"
	"Skripsigma-BE/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	tagService *service.TagService
}

func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{tagService}
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	var req dto.CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tag, err := h.tagService.CreateTag(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tag created successfully",
		"tag":     tag,
	})
}

func (h *TagHandler) GetTagByID(c *fiber.Ctx) error {
	
	id := c.Params("id")
	tag, err := h.tagService.GetTagByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tag not found"})
	}

	return c.JSON(fiber.Map{"tag": tag})
}

func (h *TagHandler) GetAllTags(c *fiber.Ctx) error {
	
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"tags": tags})
}