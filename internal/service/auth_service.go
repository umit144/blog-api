package service

import (
	"fmt"
	"go-blog/internal/database"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"os"
	"strconv"
	"time"

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

	idInt, _ := strconv.Atoi(id)
	user, err := h.userRepository.FindById(idInt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *AuthService) Login(email string, password string) (*string, error) {
	user, err := h.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(user.ID),
		"iss": "go-blog",
		"aud": "user-role", // TODO : implement authorization
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
