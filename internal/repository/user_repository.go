package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-blog/internal/types"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindAll() ([]types.User, error)
	FindByEmail(email string) (*types.User, error)
	FindById(id string) (*types.User, error)
	Create(user types.User) (*types.User, error)
	Update(id string, user types.User) (*types.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo userRepository) FindAll() ([]types.User, error) {
	var users []types.User

	sql, args, err := sq.Select("id, name, lastname, email, password, created_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindAll: %v", err)
	}

	rows, err := repo.db.QueryContext(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing FindAll query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row in FindAll: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows in FindAll: %v", err)
	}

	return users, nil
}

func (repo userRepository) FindByEmail(email string) (*types.User, error) {
	var user types.User

	sql, args, err := sq.Select("id, name, lastname, email, password, created_at").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindByEmail: %v", err)
	}

	err = repo.db.QueryRowContext(context.Background(), sql, args...).
		Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo userRepository) FindById(id string) (*types.User, error) {
	var user types.User

	sql, args, err := sq.Select("id, name, lastname, email, created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindById: %v", err)
	}

	err = repo.db.QueryRowContext(context.Background(), sql, args...).
		Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error executing FindById query: %v", err)
	}

	return &user, nil
}

func (repo userRepository) Create(user types.User) (*types.User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user.Password = string(bytes)

	sql, args, err := sq.Insert("users").
		Columns("name", "lastname", "email", "password").
		Values(user.Name, user.Lastname, user.Email, user.Password).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error creating SQL for Create: %v", err)
	}

	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fmt.Errorf("user with email %s already exists", user.Email)
		}
		return nil, fmt.Errorf("error executing Create query: %v", err)
	}

	return &user, nil
}

func (repo userRepository) Update(id string, user types.User) (*types.User, error) {
	checkSQL, checkArgs, err := sq.Select("id").
		From("users").
		Where(sq.And{
			sq.Eq{"email": user.Email},
			sq.NotEq{"id": id},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for email check: %v", err)
	}

	var existingID string
	err = repo.db.QueryRowContext(context.Background(), checkSQL, checkArgs...).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing email: %v", err)
	}
	if existingID != "" {
		return nil, fmt.Errorf("a user with email %s already exists", user.Email)
	}

	sql, args, err := sq.Update("users").
		Set("name", user.Name).
		Set("lastname", user.Lastname).
		Set("email", user.Email).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for Update: %v", err)
	}

	result, err := repo.db.ExecContext(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing Update query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no user found with id %s to update", id)
	}

	return &user, nil
}

func (repo userRepository) Delete(id string) error {
	sql, args, err := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error creating SQL for Delete: %v", err)
	}

	result, err := repo.db.ExecContext(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("error executing Delete query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %s to delete", id)
	}

	return nil
}
