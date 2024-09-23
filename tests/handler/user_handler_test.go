package handler_test

import (
	"encoding/json"
	"errors"
	"go-blog/internal/handler"
	"go-blog/internal/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]types.User, error) {
	args := m.Called()
	return args.Get(0).([]types.User), args.Error(1)
}

func (m *MockUserRepository) FindById(id string) (*types.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) Create(user types.User) (*types.User, error) {
	args := m.Called(user)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) Update(id string, user types.User) (*types.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetUserHandler(t *testing.T) {
	t.Run("Get all users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		handler := handler.NewUserHandler(mockRepo)

		users := []types.User{
			{Id: "1", Name: "John", Lastname: "Doe", Email: "john@example.com"},
			{Id: "2", Name: "Jane", Lastname: "Doe", Email: "jane@example.com"},
		}

		mockRepo.On("FindAll").Return(users, nil)

		app := fiber.New()
		app.Get("/users", handler.GetUserHandler)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []types.User
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, users, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get user by ID", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		handler := handler.NewUserHandler(mockRepo)

		user := &types.User{Id: "1", Name: "John", Lastname: "Doe", Email: "john@example.com"}

		mockRepo.On("FindById", "1").Return(user, nil)

		app := fiber.New()
		app.Get("/users/:id", handler.GetUserHandler)

		req := httptest.NewRequest("GET", "/users/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result types.User
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, *user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		handler := handler.NewUserHandler(mockRepo)

		mockRepo.On("FindById", "999").Return(nil, errors.New("user not found"))

		app := fiber.New()
		app.Get("/users/:id", handler.GetUserHandler)

		req := httptest.NewRequest("GET", "/users/999", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "User not found", result["error"])
		assert.Contains(t, result["message"], "Error retrieving user with ID 999")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error fetching all users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		handler := handler.NewUserHandler(mockRepo)

		mockRepo.On("FindAll").Return([]types.User{}, errors.New("database error"))

		app := fiber.New()
		app.Get("/users", handler.GetUserHandler)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Failed to retrieve users", result["error"])
		assert.Contains(t, result["message"], "Error listing all users")
		mockRepo.AssertExpectations(t)
	})
}
