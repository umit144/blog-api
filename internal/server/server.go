package server

import (
	"github.com/gofiber/fiber/v2"

	"go-blog/internal/database"
	"go-blog/internal/handler"
)

type FiberServer struct {
	*fiber.App

	db          database.Service
	userHandler handler.UserHandler
}

func New() *FiberServer {
	var db = database.New()
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-blog",
			AppName:      "go-blog",
		}),

		db:          db,
		userHandler: handler.NewUserHandler(db),
	}

	return server
}
