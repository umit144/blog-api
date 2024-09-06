package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	s.App.Get("/websocket", websocket.New(s.websocketHandler))

	var api = s.App.Group("/api")
	var authMiddleware = s.NewAuthMiddleware()

	var userResource = api.Group("/user")
	userResource.Get("/", authMiddleware, s.userHandler.GetUserHandler)
	userResource.Get("/:id", s.userHandler.GetUserHandler)
	userResource.Post("/", s.userHandler.CreateUserHandler)
	userResource.Put("/:id", s.userHandler.UpdateUserHandler)
	userResource.Delete("/:id", s.userHandler.DeleteUserHandler)

	var authResource = api.Group("/auth")
	authResource.Post("/login", s.authHandler.LoginHandler)
	authResource.Post("/register", s.authHandler.RegisterHandler)
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
