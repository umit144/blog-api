package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *FiberServer) RegisterFiberMiddlewares() {
	s.App.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
}

func (s *FiberServer) NewAuthMiddleware() fiber.Handler {
	return keyauth.New(keyauth.Config{
		//KeyLookup: "cookie:access_token",
		Validator: func(c *fiber.Ctx, token string) (bool, error) {
			user, err := s.authService.ParseToken(token)
			if err != nil {
				return false, keyauth.ErrMissingOrMalformedAPIKey
			}
			c.Locals("user", user)
			return true, nil
		},
	})
}
