package handler

import (
	"fmt"
	"go-blog/internal/repository"
	"go-blog/internal/types"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userRepository repository.UserRepository
}

type UserHandler interface {
	GetUserHandler(c *fiber.Ctx) error
	CreateUserHandler(c *fiber.Ctx) error
	UpdateUserHandler(c *fiber.Ctx) error
	DeleteUserHandler(c *fiber.Ctx) error
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

func (h *userHandler) CreateUserHandler(c *fiber.Ctx) error {
	var user types.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid data format",
			"message": fmt.Sprintf("Error parsing user data: %v", err),
		})
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation error",
			"fails": err,
		})
	}

	createdUser, err := h.userRepository.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create user",
			"message": fmt.Sprintf("Error occurred while creating new user: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *userHandler) UpdateUserHandler(c *fiber.Ctx) error {
	var user types.User
	id := c.Params("id")

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid data format",
			"message": fmt.Sprintf("Error parsing user data: %v", err),
		})
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation error",
			"fails": err,
		})
	}

	updatedUser, err := h.userRepository.Update(id, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update user",
			"message": fmt.Sprintf("Error updating user with ID %s: %v", id, err),
		})
	}

	return c.JSON(updatedUser)
}

func (h *userHandler) DeleteUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var err = h.userRepository.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete user",
			"message": fmt.Sprintf("Error deleting user with ID %s: %v", id, err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("User with ID %s successfully deleted", id),
	})
}
