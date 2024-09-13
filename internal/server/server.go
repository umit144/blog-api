package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-blog/internal/database"
	"go-blog/internal/handler"
	"go-blog/internal/service"
)

type FiberServer struct {
	*fiber.App

	db              database.Service
	userHandler     handler.UserHandler
	authHandler     handler.AuthHandler
	authService     service.AuthService
	postHandler     handler.PostHandler
	categoryHandler handler.CategoryHandler
}

func New() *FiberServer {
	var db = database.New()
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-blog",
			AppName:      "go-blog",
		}),

		db:              db,
		userHandler:     handler.NewUserHandler(db),
		authHandler:     handler.NewAuthHandler(db),
		authService:     service.NewAuthService(db),
		postHandler:     handler.NewPostHandler(db),
		categoryHandler: handler.NewCategoryHandler(db),
	}

	server.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${url}\n",
	}))

	server.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	return server
}
