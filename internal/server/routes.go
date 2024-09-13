package server

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"go-blog/internal/types"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	api := s.App.Group("/api")
	authMiddleware := keyauth.New(keyauth.Config{
		KeyLookup: "cookie:access_token",
		Validator: func(c *fiber.Ctx, token string) (bool, error) {
			user, err := s.authService.ParseToken(token)
			if err != nil {
				return false, keyauth.ErrMissingOrMalformedAPIKey
			}
			c.Locals("user", *user)
			return true, nil
		},
	})

	api.Get("/health", s.healthHandler)
	api.Get("/websocket", websocket.New(s.websocketHandler))

	userRoutes := api.Group("/user")
	userRoutes.Use(authMiddleware)
	{
		userRoutes.Get("/", s.userHandler.GetUserHandler)
		userRoutes.Get("/:id", s.userHandler.GetUserHandler)
		userRoutes.Post("/", s.userHandler.CreateUserHandler)
		userRoutes.Put("/:id", s.userHandler.UpdateUserHandler)
		userRoutes.Delete("/:id", s.userHandler.DeleteUserHandler)
	}

	postRoutes := api.Group("/post")
	postRoutes.Use(authMiddleware)
	{
		postRoutes.Get("/", s.postHandler.GetPostHandler)
		postRoutes.Get("/:slugOrId", s.postHandler.GetPostHandler)
		postRoutes.Post("/", s.postHandler.CreatePostHandler)
		postRoutes.Put("/:id", s.postHandler.UpdatePostHandler)
		postRoutes.Delete("/:id", s.postHandler.DeletePostHandler)
	}

	authRoutes := api.Group("/auth")
	{
		authRoutes.Post("/login", s.authHandler.LoginHandler)
		authRoutes.Post("/register", s.authHandler.RegisterHandler)
		authRoutes.Get("/google/login", s.authHandler.GoogleLoginHandler)
		authRoutes.Post("/google/callback", s.authHandler.GoogleCallbackHandler)
	}

	authRoutes.Get("/session", authMiddleware, func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(types.User)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Session not found",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(user)
	})
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
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
