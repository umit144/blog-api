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

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo UserRepository) FindAll() ([]types.User, error) {
	var users []types.User

	sql, _, _ := sq.Select("id, name, lastname, email, password, created_at").From("users").ToSql()
	rows, err := repo.db.QueryContext(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo UserRepository) FindByEmail(email string) (*types.User, error) {
	var user types.User

	sql, args, err := sq.Select("id, name, lastname, email, password, created_at").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRowContext(context.Background(), sql, args...)
	err = row.Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo UserRepository) FindById(id string) (*types.User, error) {
	var user types.User

	sql, args, err := sq.Select("id, name, lastname, email, created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRowContext(context.Background(), sql, args...)
	err = row.Scan(&user.Id, &user.Name, &user.Lastname, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo UserRepository) Create(user types.User) (*types.User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, err
	}

	user.Password = string(bytes)

	sql, args, err := sq.Insert("users").
		Columns("name", "lastname", "email", "password").
		Values(user.Name, user.Lastname, user.Email, user.Password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(sql, args...).Scan(&user.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fmt.Errorf("user with this email already exists")
		}
		return nil, err
	}

	return &user, nil
}

func (repo UserRepository) Update(id int, user types.User) (*types.User, error) {
	sql, args, err := sq.Update("users").
		Set("name", user.Name).
		Set("lastname", user.Lastname).
		Set("email", user.Email).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	result, err := repo.db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no rows affected, user with id %d not found", id)
	}

	return &user, nil
}

func (repo UserRepository) Delete(id int) error {
	sql, args, err := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	result, err := repo.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, user with id %d not found", id)
	}

	return nil
}
