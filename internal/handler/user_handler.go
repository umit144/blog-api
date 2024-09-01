package handler

import (
	"go-blog/internal/database"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	HandleGet(c *fiber.Ctx) error
	HandlePost(c *fiber.Ctx) error
	HandlePut(c *fiber.Ctx) error
	HandlePatch(c *fiber.Ctx) error
	HandleDelete(c *fiber.Ctx) error
}

type postgresUserHandler struct {
	userRepository repository.UserRepository
}

func NewUserHandler(db database.Service) UserHandler {
	return &postgresUserHandler{
		userRepository: *repository.NewUserRepository(db.GetInstance()),
	}
}

func (h *postgresUserHandler) HandleGet(c *fiber.Ctx) error {
	var id = c.Params("id")

	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Can't parse id to int",
				"message": err.Error(),
			})
		}
		user, err := h.userRepository.FindById(idInt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to getting user by id",
				"message": err.Error(),
			})
		}
		return c.JSON(user)
	}

	var users, err = h.userRepository.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to getting users",
			"message": err.Error(),
		})
	}

	return c.JSON(users)
}

func (h *postgresUserHandler) HandlePost(c *fiber.Ctx) error {
	var user types.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Can't parse payload",
			"message": err.Error(),
		})
	}

	if err := user.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	createdUser, err := h.userRepository.Create(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create user",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(createdUser)
}

func (h *postgresUserHandler) HandlePut(c *fiber.Ctx) error {
	var user types.User
	var id = c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Can't parse id to int",
			"message": err.Error(),
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Can't parse payload",
			"message": err.Error(),
		})
	}

	if err := user.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	updatedUser, err := h.userRepository.Update(idInt, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update user",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(updatedUser)
}

func (h *postgresUserHandler) HandlePatch(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}

func (h *postgresUserHandler) HandleDelete(c *fiber.Ctx) error {
	var id = c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Can't parse id to int",
			"message": err.Error(),
		})
	}

	err = h.userRepository.Delete(idInt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete user",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
