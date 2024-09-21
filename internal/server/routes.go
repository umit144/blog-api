package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/openapi.yml",
		Path:     "swagger",
		Title:    "Go Blog API Docs",
	}))

	api := s.App.Group("/api")
	authMiddleware := keyauth.New(keyauth.Config{
		KeyLookup:    "cookie:access_token",
		Validator:    s.authService.ValidateSession,
		ErrorHandler: s.authHandler.AuthFailHandler,
	})

	api.Get("/health", s.healthHandler)
	api.Get("/websocket", websocket.New(s.websocketHandler))

	authRoutes := api.Group("/auth")
	{
		authRoutes.Post("/login", s.authHandler.LoginHandler)
		authRoutes.Post("/register", s.authHandler.RegisterHandler)
		authRoutes.Get("/session", authMiddleware, s.authHandler.SessionHandler)
		authRoutes.Get("/google/login", s.authHandler.GoogleLoginHandler)
		authRoutes.Post("/google/callback", s.authHandler.GoogleCallbackHandler)
		authRoutes.Get("/logout", s.authHandler.LogoutHandler)
	}

	userRoutes := api.Group("/users")
	userRoutes.Use(authMiddleware)
	{
		userRoutes.Get("/", s.userHandler.GetUserHandler)
		userRoutes.Get("/:id", s.userHandler.GetUserHandler)
		userRoutes.Post("/", s.userHandler.CreateUserHandler)
		userRoutes.Put("/:id", s.userHandler.UpdateUserHandler)
		userRoutes.Delete("/:id", s.userHandler.DeleteUserHandler)
	}

	postRoutes := api.Group("/posts")
	{
		postRoutes.Get("/", s.postHandler.GetPostHandler)
		postRoutes.Get("/:slugOrId", s.postHandler.GetPostHandler)
		postRoutes.Post("/", authMiddleware, s.postHandler.CreatePostHandler)
		postRoutes.Put("/:id", authMiddleware, s.postHandler.UpdatePostHandler)
		postRoutes.Delete("/:id", authMiddleware, s.postHandler.DeletePostHandler)
		postRoutes.Post("/:postId/categories/:categoryId", authMiddleware, s.postHandler.AssignCategoryToPostHandler)
		postRoutes.Delete("/:postId/categories/:categoryId", authMiddleware, s.postHandler.UnassignCategoryFromPostHandler)
		postRoutes.Get("/:postId/categories", s.postHandler.GetCategoriesForPostHandler)
		postRoutes.Put("/:postId/categories", authMiddleware, s.postHandler.UpdatePostCategoriesHandler)
	}

	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.Get("/", s.categoryHandler.GetCategoryHandler)
		categoryRoutes.Get("/:slugOrId", s.categoryHandler.GetCategoryHandler)
		categoryRoutes.Post("/", authMiddleware, s.categoryHandler.CreateCategoryHandler)
		categoryRoutes.Put("/:id", authMiddleware, s.categoryHandler.UpdateCategoryHandler)
		categoryRoutes.Delete("/:id", authMiddleware, s.categoryHandler.DeleteCategoryHandler)
	}

	fileRoutes := api.Group("/files")
	fileRoutes.Use(authMiddleware)
	{
		fileRoutes.Post("/", s.fileHandler.UploadFileHandler)
		fileRoutes.Delete("/:filename", s.fileHandler.DeleteFileHandler)
	}
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.dbStatus)
}

func (s *FiberServer) websocketHandler(con *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			_, _, err := con.ReadMessage()
			if err != nil {
				cancel()
				log.Println("Receiver Closing", err)
				break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
			if err := con.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
				log.Printf("could not write to socket: %v", err)
				return
			}
			time.Sleep(time.Second * 2)
		}
	}
}
