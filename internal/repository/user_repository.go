package repository

import (
	"context"
	"database/sql"
	"go-blog/internal/types"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo UserRepository) FindAll() ([]types.User, error) {
	var users []types.User

	sql, _, _ := sq.Select("id, name, lastname, email, created_at").From("users").ToSql()
	rows, err := repo.db.QueryContext(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Name, &user.Lastname, &user.Email, &user.CreatedAt)
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

func (repo UserRepository) FindById(id int) (*types.User, error) {
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
	err = row.Scan(&user.ID, &user.Name, &user.Lastname, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo UserRepository) Create(user types.User) (*types.User, *error) {

	return nil, nil
}

func (repo UserRepository) Update(id int, user types.User) (*types.User, *error) {
	return nil, nil
}

func (repo UserRepository) Delete(id int) *error {
	return nil
}
