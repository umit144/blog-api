package repository_test

import (
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewCategoryRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Category 1", "category-1", time.Now()).
		AddRow("2", "Category 2", "category-2", time.Now())

	mock.ExpectQuery("SELECT id, title, slug, created_at FROM categories").WillReturnRows(rows)

	categories, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Category 1", categories[0].Title)
	assert.Equal(t, "category-2", categories[1].Slug)
}

func TestCategoryRepository_FindBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewCategoryRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Category 1", "category-1", time.Now())

	mock.ExpectQuery("SELECT id, title, slug, created_at FROM categories WHERE").
		WithArgs("category-1").
		WillReturnRows(rows)

	category, err := repo.FindBySlug("category-1")

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Category 1", category.Title)
}

func TestCategoryRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewCategoryRepository(db)

	mock.ExpectQuery("SELECT EXISTS").WithArgs("new-category").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "New Category", "new-category", time.Now())

	mock.ExpectQuery("INSERT INTO categories").
		WithArgs("New Category", "new-category").
		WillReturnRows(rows)

	newCategory := types.Category{
		Title: "New Category",
		Slug:  "new-category",
	}

	createdCategory, err := repo.Create(newCategory)

	assert.NoError(t, err)
	assert.NotNil(t, createdCategory)
	assert.Equal(t, "New Category", createdCategory.Title)
	assert.Equal(t, "new-category", createdCategory.Slug)
}

func TestCategoryRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewCategoryRepository(db)

	mock.ExpectQuery("SELECT EXISTS").WithArgs("updated-category", "1").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Updated Category", "updated-category", time.Now())

	mock.ExpectQuery("UPDATE categories").
		WithArgs("Updated Category", "updated-category", "1").
		WillReturnRows(rows)

	updatedCategory := types.Category{
		Title: "Updated Category",
		Slug:  "updated-category",
	}

	result, err := repo.Update("1", updatedCategory)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Category", result.Title)
	assert.Equal(t, "updated-category", result.Slug)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewCategoryRepository(db)

	mock.ExpectExec("DELETE FROM categories").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete("1")

	assert.NoError(t, err)
}

func TestCategoryRepository_FindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewCategoryRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Category 1", "category-1", time.Now())

	mock.ExpectQuery("SELECT id, title, slug, created_at FROM categories WHERE").
		WithArgs("1").
		WillReturnRows(rows)

	category, err := repo.FindById("1")

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Category 1", category.Title)
	assert.Equal(t, "category-1", category.Slug)
}
