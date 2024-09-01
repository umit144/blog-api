package server

import (
	"github.com/gofiber/fiber/v2"

	"go-blog/internal/database"
	"go-blog/internal/handler"
	"go-blog/internal/service"
)

type FiberServer struct {
	*fiber.App

	db          database.Service
	userHandler handler.UserHandler
	authHandler handler.AuthHandler
	authService service.AuthService
}

func New() *FiberServer {
	var db = database.New()
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-blog",
			AppName:      "go-blog",
		}),

		db:          db,
		userHandler: *handler.NewUserHandler(db),
		authHandler: *handler.NewAuthHandler(db),
		authService: *service.NewAuthService(db),
	}

	return server
}
