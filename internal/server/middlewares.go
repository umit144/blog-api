package server

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *FiberServer) RegisterFiberMiddlewares() {
	s.App.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
}
