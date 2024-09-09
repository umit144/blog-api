package service

import (
	"database/sql"
	"errors"
	"fmt"
	"go-blog/internal/database"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"os"
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

func (s *AuthService) ParseToken(tokenString string) (*types.User, error) {
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

func (s *AuthService) Login(email string, password string) (*string, *types.User, error) {
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

func (s *AuthService) Register(user types.User) (*string, *types.User, error) {
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

func (s *AuthService) LoginOrRegisterWithGoogle(email, name, googleID, profilePicture string) (*string, *types.User, error) {
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
