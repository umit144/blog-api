package service_test

import (
	"database/sql"
	"errors"
	"go-blog/internal/service"
	"go-blog/internal/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]types.User, error) {
	args := m.Called()
	return args.Get(0).([]types.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) FindById(id string) (*types.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) Create(user types.User) (*types.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) Update(id string, user types.User) (*types.User, error) {
	args := m.Called(id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func TestRegister(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	testUser := types.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("Successful registration", func(t *testing.T) {
		createdUser := testUser
		createdUser.Id = "123"
		createdUser.CreatedAt = time.Now()
		createdUser.UpdatedAt = time.Now()

		mockRepo.On("Create", mock.AnythingOfType("types.User")).Return(&createdUser, nil).Once()

		token, user, err := authService.Register(testUser)

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, createdUser.Email, user.Email)
		assert.Equal(t, createdUser.Id, user.Id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Registration failure", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("types.User")).Return(nil, errors.New("registration failed")).Once()

		token, user, err := authService.Register(testUser)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "registration failed")
		mockRepo.AssertExpectations(t)
	})
}
func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	// Create a bcrypt hash of the password "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	testCases := []struct {
		name          string
		email         string
		password      string
		mockUser      *types.User
		mockError     error
		expectedError bool
	}{
		{
			name:     "Successful login",
			email:    "test@example.com",
			password: "password123",
			mockUser: &types.User{
				Id:       "1",
				Email:    "test@example.com",
				Password: string(hashedPassword),
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "User not found",
			email:         "nonexistent@example.com",
			password:      "password123",
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: true,
		},
		{
			name:     "Incorrect password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockUser: &types.User{
				Id:       "1",
				Email:    "test@example.com",
				Password: string(hashedPassword),
			},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("FindByEmail", tc.email).Return(tc.mockUser, tc.mockError).Once()

			token, user, err := authService.Login(tc.email, tc.password)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotNil(t, user)
				assert.Equal(t, tc.mockUser.Id, user.Id)
				assert.Equal(t, tc.mockUser.Email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestLoginOrRegisterWithGoogle(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	testCases := []struct {
		name           string
		email          string
		findByEmailErr error
		createUserErr  error
		updateUserErr  error
		expectedError  bool
	}{
		{
			name:           "Existing user login",
			email:          "existing@example.com",
			findByEmailErr: nil,
			updateUserErr:  nil,
			expectedError:  false,
		},
		{
			name:           "New user registration",
			email:          "new@example.com",
			findByEmailErr: sql.ErrNoRows,
			createUserErr:  nil,
			expectedError:  false,
		},
		{
			name:           "Error finding user",
			email:          "error@example.com",
			findByEmailErr: errors.New("database error"),
			expectedError:  true,
		},
		{
			name:           "Error creating new user",
			email:          "newerror@example.com",
			findByEmailErr: sql.ErrNoRows,
			createUserErr:  errors.New("creation error"),
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUser := &types.User{
				Id:           "123",
				Email:        tc.email,
				Name:         "Test User",
				GoogleID:     "123456789",
				AuthProvider: "google",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if tc.findByEmailErr == nil {
				mockRepo.On("FindByEmail", tc.email).Return(mockUser, nil).Once()
				mockRepo.On("Update", mockUser.Id, mock.AnythingOfType("types.User")).Return(mockUser, tc.updateUserErr).Once()
			} else {
				mockRepo.On("FindByEmail", tc.email).Return(nil, tc.findByEmailErr).Once()
				if tc.findByEmailErr == sql.ErrNoRows {
					if tc.createUserErr == nil {
						mockRepo.On("Create", mock.AnythingOfType("types.User")).Return(mockUser, nil).Once()
					} else {
						mockRepo.On("Create", mock.AnythingOfType("types.User")).Return(nil, tc.createUserErr).Once()
					}
				}
			}

			token, user, err := authService.LoginOrRegisterWithGoogle(tc.email, "Test User", "123456789", "")

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotNil(t, user)
				assert.Equal(t, tc.email, user.Email)
				assert.Equal(t, "123456789", user.GoogleID)
				assert.Equal(t, "google", user.AuthProvider)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
