package repository_test

import (
	"database/sql"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/lib/pq"
	"go-blog/internal/repository"
	"go-blog/internal/types"
)

func TestUserRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password", "created_at"}).
		AddRow("1", "John", "Doe", "john@example.com", "hashedpassword1", time.Now()).
		AddRow("2", "Jane", "Doe", "jane@example.com", "hashedpassword2", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	users, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "John", users[0].Name)
	assert.Equal(t, "Jane", users[1].Name)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "password", "created_at"}).
		AddRow("1", "John", "Doe", "john@example.com", "hashedpassword", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("john@example.com").WillReturnRows(rows)

	user, err := repo.FindByEmail("john@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestUserRepository_FindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "email", "created_at"}).
		AddRow("1", "John", "Doe", "john@example.com", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("1").WillReturnRows(rows)

	user, err := repo.FindById("1")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "1", user.Id)
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	newUser := types.User{
		Name:     "John",
		Lastname: "Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(newUser.Name, newUser.Lastname, newUser.Email, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("1", time.Now()))

	createdUser, err := repo.Create(newUser)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "John", createdUser.Name)
	assert.Equal(t, "1", createdUser.Id)

	err = bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("password123"))
	assert.NoError(t, err)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	userID := uuid.New().String()
	updatedUser := types.User{
		Id:       userID,
		Name:     "John Updated",
		Lastname: "Doe Updated",
		Email:    "john.updated@example.com",
	}

	mock.ExpectQuery("SELECT id FROM users WHERE").
		WithArgs(updatedUser.Email, userID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("UPDATE users").
		WithArgs(updatedUser.Name, updatedUser.Lastname, updatedUser.Email, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := repo.Update(userID, updatedUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John Updated", result.Name)
	assert.Equal(t, userID, result.Id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	userID := uuid.New().String()

	mock.ExpectExec("DELETE FROM users").WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(userID)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)

	newUser := types.User{
		Name:     "John",
		Lastname: "Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	pgError := &pq.Error{
		Code:    "23505",
		Message: "duplicate key value violates unique constraint",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(newUser.Name, newUser.Lastname, newUser.Email, sqlmock.AnyArg()).
		WillReturnError(pgError)

	_, err = repo.Create(newUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user with email john@example.com already exists")
}
