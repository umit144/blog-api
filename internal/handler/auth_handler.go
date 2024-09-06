package handler

import (
	"go-blog/internal/database"
	"go-blog/internal/service"
	"go-blog/internal/types"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(db database.Service) *AuthHandler {
	return &AuthHandler{
		authService: *service.NewAuthService(db),
	}
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Can't parse payload",
			"message": err.Error(),
		})
	}

	token, user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to login",
			"message": err.Error(),
		})
	}

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{}

	authenticatedUser.Token = *token
	authenticatedUser.User = *user

	return c.JSON(authenticatedUser)
}

func (h *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
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

	token, createdUser, err := h.authService.Register(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create user",
			"message": err.Error(),
		})
	}

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{}

	authenticatedUser.Token = *token
	authenticatedUser.User = *createdUser

	return c.JSON(authenticatedUser)
}
