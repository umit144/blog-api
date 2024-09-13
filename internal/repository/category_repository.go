package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-blog/internal/types"

	sq "github.com/Masterminds/squirrel"
)

type CategoryRepository interface {
	FindAll() ([]types.Category, error)
	FindBySlug(slug string) (*types.Category, error)
	FindById(id string) (*types.Category, error)
	Create(category types.Category) (*types.Category, error)
	Update(id string, category types.Category) (*types.Category, error)
	Delete(id string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (repo categoryRepository) FindAll() ([]types.Category, error) {
	var categories []types.Category

	sql, args, err := sq.Select("id, title, slug, created_at").
		From("categories").
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
		var category types.Category
		err := rows.Scan(
			&category.Id,
			&category.Title,
			&category.Slug,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row in FindAll: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows in FindAll: %v", err)
	}

	return categories, nil
}

func (repo categoryRepository) FindBySlug(slug string) (*types.Category, error) {
	query := sq.Select("id, title, slug, created_at").
		From("categories").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"slug": slug})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindBySlug: %v", err)
	}

	var category types.Category
	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&category.Id,
		&category.Title,
		&category.Slug,
		&category.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing FindBySlug query: %v", err)
	}

	return &category, nil
}

func (repo categoryRepository) FindById(id string) (*types.Category, error) {
	query := sq.Select("id, title, slug, created_at").
		From("categories").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindById: %v", err)
	}

	var category types.Category
	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&category.Id,
		&category.Title,
		&category.Slug,
		&category.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing FindById query: %v", err)
	}

	return &category, nil
}

func (repo categoryRepository) Create(category types.Category) (*types.Category, error) {
	slug, err := repo.generateUniqueSlug(category.Slug, "")
	if err != nil {
		return nil, fmt.Errorf("error generating unique slug: %v", err)
	}

	insertQuery := sq.Insert("categories").
		Columns("title", "slug").
		Values(category.Title, slug).
		Suffix("RETURNING id, title, slug, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for Create: %v", err)
	}

	var createdCategory types.Category

	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&createdCategory.Id,
		&createdCategory.Title,
		&createdCategory.Slug,
		&createdCategory.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing Create query: %v", err)
	}

	return &createdCategory, nil
}

func (repo categoryRepository) Update(id string, category types.Category) (*types.Category, error) {
	slug, err := repo.generateUniqueSlug(category.Slug, id)
	if err != nil {
		return nil, fmt.Errorf("error generating unique slug for update: %v", err)
	}

	updateQuery := sq.Update("categories").
		Set("title", category.Title).
		Set("slug", slug).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, title, slug, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := updateQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for Update: %v", err)
	}

	var updatedCategory types.Category
	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&updatedCategory.Id,
		&updatedCategory.Title,
		&updatedCategory.Slug,
		&updatedCategory.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing Update query: %v", err)
	}

	return &updatedCategory, nil
}

func (repo categoryRepository) Delete(id string) error {
	deleteQuery := sq.Delete("categories").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := deleteQuery.ToSql()
	if err != nil {
		return fmt.Errorf("error creating SQL for Delete: %v", err)
	}

	result, err := repo.db.ExecContext(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("error executing Delete query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no category found with id %s to delete", id)
	}

	return nil
}

func (repo categoryRepository) generateUniqueSlug(baseSlug string, excludeId string) (string, error) {
	slug := baseSlug
	counter := 1

	for {
		var exists bool
		var err error

		if excludeId != "" {
			err = repo.db.QueryRowContext(
				context.Background(),
				"SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1 AND id != $2)",
				slug, excludeId,
			).Scan(&exists)
		} else {
			err = repo.db.QueryRowContext(
				context.Background(),
				"SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1)",
				slug,
			).Scan(&exists)
		}

		if err != nil {
			return "", fmt.Errorf("error checking slug existence: %v", err)
		}

		if !exists {
			break
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug, nil
}
