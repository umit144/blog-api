package service

import (
	"fmt"
	"go-blog/internal/database"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte = []byte(os.Getenv("JWT_SECRET"))

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(db database.Service) *AuthService {
	return &AuthService{
		userRepository: *repository.NewUserRepository(db.GetInstance()),
	}
}

func (h *AuthService) ParseToken(tokenString string) (*types.User, error) {
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

	user, err := h.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *AuthService) Login(email string, password string) (*string, *types.User, error) {
	user, err := h.userRepository.FindByEmail(email)
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

func (h *AuthService) Register(user types.User) (*string, *types.User, error) {
	createdUser, err := h.userRepository.Create(user)
	if err != nil {
		return nil, nil, err
	}

	token, err := createdUser.CreateToken()
	if err != nil {
		return nil, nil, err
	}

	return token, createdUser, nil
}
