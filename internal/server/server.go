package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-blog/internal/database"
	"go-blog/internal/handler"
	"go-blog/internal/repository"
	"go-blog/internal/service"
)

type FiberServer struct {
	*fiber.App

	dbStatus        map[string]string
	userHandler     handler.UserHandler
	authHandler     handler.AuthHandler
	postHandler     handler.PostHandler
	categoryHandler handler.CategoryHandler
	fileHandler     handler.FileHandler

	authService service.AuthService
}

func New() *FiberServer {
	var db = database.New()

	var userRepository = repository.NewUserRepository(db.GetInstance())
	var postRepository = repository.NewPostRepository(db.GetInstance())
	var categoryRepository = repository.NewCategoryRepository(db.GetInstance())

	var authService = service.NewAuthService(userRepository)
	var fileService = service.NewFileService()

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-blog",
			AppName:      "go-blog",
		}),
		dbStatus:        db.Health(),
		userHandler:     handler.NewUserHandler(userRepository),
		authHandler:     handler.NewAuthHandler(authService),
		postHandler:     handler.NewPostHandler(postRepository),
		categoryHandler: handler.NewCategoryHandler(categoryRepository),
		fileHandler:     handler.NewFileHandler(fileService),
		authService:     authService,
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
