package server

import (
	"github.com/gofiber/fiber/v2"

	"go-blog/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-blog",
			AppName:      "go-blog",
		}),

		db: database.New(),
	}

	return server
}
