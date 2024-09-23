package handler

import (
	"fmt"
	"go-blog/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userRepository repository.UserRepository
}

type UserHandler interface {
	GetUserHandler(c *fiber.Ctx) error
}

func NewUserHandler(userRepository repository.UserRepository) UserHandler {
	return &userHandler{userRepository}
}

func (h *userHandler) GetUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id != "" {
		user, err := h.userRepository.FindById(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "User not found",
				"message": fmt.Sprintf("Error retrieving user with ID %s: %v", id, err),
			})
		}
		return c.JSON(user)
	}

	users, err := h.userRepository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve users",
			"message": fmt.Sprintf("Error listing all users: %v", err),
		})
	}

	return c.JSON(users)
}
