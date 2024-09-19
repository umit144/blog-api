package service

import (
	"database/sql"
	"errors"
	"fmt"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte = []byte(os.Getenv("JWT_SECRET"))

type authService struct {
	userRepository repository.UserRepository
}

type AuthService interface {
	ParseToken(tokenString string) (*types.User, error)
	Login(email string, password string) (*string, *types.User, error)
	Register(user types.User) (*string, *types.User, error)
	LoginOrRegisterWithGoogle(email, name, googleID, profilePicture string) (*string, *types.User, error)
	GenerateAuthCookie(token string) *fiber.Cookie
	ValidateSession(c *fiber.Ctx, token string) (bool, error)
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepository}
}

func (s *authService) ParseToken(tokenString string) (*types.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email string, password string) (*string, *types.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, err
	}

	tokenString, err := user.CreateToken()
	if err != nil {
		return nil, nil, err
	}

	return tokenString, user, nil
}

func (s *authService) Register(user types.User) (*string, *types.User, error) {
	createdUser, err := s.userRepository.Create(user)
	if err != nil {
		return nil, nil, err
	}

	token, err := createdUser.CreateToken()
	if err != nil {
		return nil, nil, err
	}

	return token, createdUser, nil
}

func (s *authService) LoginOrRegisterWithGoogle(email, name, googleID, profilePicture string) (*string, *types.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newUser := types.User{
				Email:          email,
				Name:           name,
				GoogleID:       googleID,
				ProfilePicture: profilePicture,
				AuthProvider:   "google",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			user, err = s.userRepository.Create(newUser)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create user: %v", err)
			}
		} else {
			return nil, nil, fmt.Errorf("error finding user: %v", err)
		}
	} else {
		user.GoogleID = googleID
		user.ProfilePicture = profilePicture
		user.AuthProvider = "google"
		user.UpdatedAt = time.Now()
		user, err = s.userRepository.Update(user.Id, *user)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update user: %v", err)
		}
	}

	token, err := user.CreateToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create token: %v", err)
	}

	return token, user, nil
}

func (s *authService) GenerateAuthCookie(token string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
}

func (s *authService) ValidateSession(c *fiber.Ctx, token string) (bool, error) {
	user, err := s.ParseToken(token)
	if err != nil {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}
	c.Locals("user", *user)
	return true, nil
}
