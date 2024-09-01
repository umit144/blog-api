package handler

import (
	"go-blog/internal/database"
	"go-blog/internal/repository"

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
	var users, err = h.userRepository.FindAll()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(users)
}

func (h *postgresUserHandler) HandlePost(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}

func (h *postgresUserHandler) HandlePut(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}

func (h *postgresUserHandler) HandlePatch(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}

func (h *postgresUserHandler) HandleDelete(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}
