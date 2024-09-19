package handler

import (
	"encoding/json"
	"fmt"
	"go-blog/internal/database"
	"go-blog/internal/service"
	"go-blog/internal/types"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler interface {
	LoginHandler(c *fiber.Ctx) error
	RegisterHandler(c *fiber.Ctx) error
	SessionHandler(c *fiber.Ctx) error
	GoogleLoginHandler(c *fiber.Ctx) error
	GoogleCallbackHandler(c *fiber.Ctx) error
	LogoutHandler(c *fiber.Ctx) error
	AuthFailHandler(c *fiber.Ctx, err error) error
}

type authHandler struct {
	authService       service.AuthService
	googleOauthConfig *oauth2.Config
}

func NewAuthHandler(db database.Service) AuthHandler {
	return &authHandler{
		authService: service.NewAuthService(db),
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("CLIENT_URL") + os.Getenv("OAUTH_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *authHandler) LoginHandler(c *fiber.Ctx) error {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing login data: %v", err),
		})
	}

	token, user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authentication failed",
			"message": fmt.Sprintf("Login attempt failed: %v", err),
		})
	}

	c.Cookie(h.authService.GenerateAuthCookie(*token))

	return c.JSON(user)
}

func (h *authHandler) RegisterHandler(c *fiber.Ctx) error {
	var user types.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing registration data: %v", err),
		})
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	token, createdUser, err := h.authService.Register(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Registration failed",
			"message": fmt.Sprintf("Error creating new user: %v", err),
		})
	}

	c.Cookie(h.authService.GenerateAuthCookie(*token))

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *authHandler) SessionHandler(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Session not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *authHandler) GoogleLoginHandler(c *fiber.Ctx) error {
	url := h.googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *authHandler) GoogleCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := h.googleOauthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Code exchange failed",
			"message": err.Error(),
		})
	}

	client := h.googleOauthConfig.Client(c.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed getting user info",
			"message": err.Error(),
		})
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		VerifiedEmail bool   `json:"verified_email"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
		HD            string `json:"hd"`
		ID            string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed decoding user info",
			"message": err.Error(),
		})
	}

	jwtToken, user, err := h.authService.LoginOrRegisterWithGoogle(googleUser.Email, googleUser.Name, googleUser.ID, googleUser.Picture)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Authentication failed",
			"message": err.Error(),
		})
	}

	c.Cookie(h.authService.GenerateAuthCookie(*jwtToken))

	return c.JSON(user)
}

func (h *authHandler) LogoutHandler(c *fiber.Ctx) error {
	c.Cookie(h.authService.GenerateAuthCookie(""))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *authHandler) AuthFailHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "Unauthorized",
		"message": err.Error(),
	})
}
