package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"

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
	postHandler handler.PostHandler
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
		postHandler: *handler.NewPostHandler(db),
	}

	server.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CLIENT_URL"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	return server
}
