package handler

import (
	"fmt"
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
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing login data: %v", err),
		})
	}

	token, user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authentication failed",
			"message": fmt.Sprintf("Login attempt failed: %v", err),
		})
	}

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{
		Token: *token,
		User:  *user,
	}

	return c.JSON(authenticatedUser)
}

func (h *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	var user types.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing registration data: %v", err),
		})
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	token, createdUser, err := h.authService.Register(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Registration failed",
			"message": fmt.Sprintf("Error creating new user: %v", err),
		})
	}

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{
		Token: *token,
		User:  *createdUser,
	}

	return c.Status(fiber.StatusCreated).JSON(authenticatedUser)
}
