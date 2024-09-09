package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-blog/internal/database"
	"go-blog/internal/service"
	"go-blog/internal/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

type AuthHandler struct {
	authService       service.AuthService
	googleOauthConfig *oauth2.Config
}

func NewAuthHandler(db database.Service) *AuthHandler {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("CLIENT_URL") + os.Getenv("OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return &AuthHandler{
		authService:       *service.NewAuthService(db),
		googleOauthConfig: googleOauthConfig,
	}
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
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

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{
		Token: *token,
		User:  *user,
	}

	return c.JSON(authenticatedUser)
}

func (h *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
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

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{
		Token: *token,
		User:  *createdUser,
	}

	return c.Status(fiber.StatusCreated).JSON(authenticatedUser)
}

func (h *AuthHandler) GoogleLoginHandler(c *fiber.Ctx) error {
	url := h.googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallbackHandler(c *fiber.Ctx) error {
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

	authenticatedUser := struct {
		Token string     `json:"token"`
		User  types.User `json:"user"`
	}{
		Token: *jwtToken,
		User:  *user,
	}

	return c.JSON(authenticatedUser)
}
