package handler

import (
	"fmt"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryHandler interface {
	GetCategoryHandler(c *fiber.Ctx) error
	CreateCategoryHandler(c *fiber.Ctx) error
	UpdateCategoryHandler(c *fiber.Ctx) error
	DeleteCategoryHandler(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryHandler(categoryRepository repository.CategoryRepository) CategoryHandler {
	return &categoryHandler{categoryRepository}
}

func (h *categoryHandler) GetCategoryHandler(c *fiber.Ctx) error {
	slugOrId := c.Params("slugOrId")

	if slugOrId != "" {
		var category *types.Category

		_, err := uuid.Parse(slugOrId)
		if err == nil {
			category, err = h.categoryRepository.FindById(slugOrId)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Category not found",
					"message": fmt.Sprintf("No category found with id: %s", slugOrId),
				})
			}
		} else {
			category, err = h.categoryRepository.FindBySlug(slugOrId)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Category not found",
					"message": fmt.Sprintf("No category found with slug: %s", slugOrId),
				})
			}
		}

		return c.JSON(category)
	}

	categories, err := h.categoryRepository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve categories",
			"message": fmt.Sprintf("Error occurred while fetching categories: %v", err),
		})
	}

	return c.JSON(categories)
}
func (h *categoryHandler) CreateCategoryHandler(c *fiber.Ctx) error {
	var category types.Category

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing category data: %v", err),
		})
	}

	if err := category.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	// Generate slug
	slug := strings.ToLower(category.Title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	category.Slug = strings.Trim(slug, "-")

	createdCategory, err := h.categoryRepository.Create(category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create category",
			"message": fmt.Sprintf("Error occurred while creating category: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdCategory)
}

func (h *categoryHandler) UpdateCategoryHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var category types.Category

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing category data: %v", err),
		})
	}

	if err := category.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	_, err := h.categoryRepository.FindById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Category not found",
			"message": fmt.Sprintf("No category found with ID: %s", id),
		})
	}

	slug := strings.ToLower(category.Title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	category.Slug = strings.Trim(slug, "-")

	updatedCategory, err := h.categoryRepository.Update(id, category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update category",
			"message": fmt.Sprintf("Error occurred while updating category: %v", err),
		})
	}

	return c.JSON(updatedCategory)
}

func (h *categoryHandler) DeleteCategoryHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := h.categoryRepository.FindById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Category not found",
			"message": fmt.Sprintf("No category found with ID: %s", id),
		})
	}

	if err := h.categoryRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete category",
			"message": fmt.Sprintf("Error occurred while deleting category: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
