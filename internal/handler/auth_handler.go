package handler

import (
	"go-blog/internal/database"
	"go-blog/internal/service"

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

	token, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to login",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
