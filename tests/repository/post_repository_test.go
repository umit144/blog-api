package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-blog/internal/repository"
	"go-blog/internal/types"
)

func TestPostRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email", "category_id", "category_title", "category_slug", "category_created_at"}).
		AddRow("1", "Test Post", "test-post", "Content", time.Now(), "1", "John", "Doe", "john@example.com", "1", "Category 1", "category-1", time.Now()).
		AddRow("1", "Test Post", "test-post", "Content", time.Now(), "1", "John", "Doe", "john@example.com", "2", "Category 2", "category-2", time.Now()).
		AddRow("2", "Another Post", "another-post", "More Content", time.Now(), "2", "Jane", "Doe", "jane@example.com", nil, nil, nil, nil)

	mock.ExpectQuery("SELECT DISTINCT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email, categories.id, categories.title, categories.slug, categories.created_at FROM posts").WillReturnRows(rows)

	posts, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Test Post", posts[0].Title)
	assert.Equal(t, "Another Post", posts[1].Title)
	assert.Len(t, posts[0].Categories, 2)
	assert.Len(t, posts[1].Categories, 0)
}

func TestPostRepository_FindAllPaginated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email", "category_id", "category_title", "category_slug", "category_created_at"}).
		AddRow("1", "Test Post", "test-post", "Content", time.Now(), "1", "John", "Doe", "john@example.com", "1", "Category 1", "category-1", time.Now()).
		AddRow("2", "Another Post", "another-post", "More Content", time.Now(), "2", "Jane", "Doe", "jane@example.com", nil, nil, nil, nil)

	mock.ExpectQuery("SELECT DISTINCT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email, categories.id, categories.title, categories.slug, categories.created_at FROM posts").WillReturnRows(rows)

	posts, totalCount, err := repo.FindAllPaginated(1, 5)

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, 10, totalCount)
}

func TestPostRepository_FindBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email"}).
		AddRow("1", "Test Post", "test-post", "Content", time.Now(), "1", "John", "Doe", "john@example.com")

	mock.ExpectQuery("SELECT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email FROM posts").WithArgs("test-post").WillReturnRows(rows)

	categoryRows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Category 1", "category-1", time.Now())

	mock.ExpectQuery("SELECT categories.id, categories.title, categories.slug, categories.created_at FROM categories").WillReturnRows(categoryRows)

	post, err := repo.FindBySlug("test-post")

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "test-post", post.Slug)
	assert.Len(t, post.Categories, 1)
}

func TestPostRepository_FindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email"}).
		AddRow("1", "Test Post", "test-post", "Content", time.Now(), "1", "John", "Doe", "john@example.com")

	mock.ExpectQuery("SELECT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email FROM posts").WithArgs("1").WillReturnRows(rows)

	categoryRows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Category 1", "category-1", time.Now())

	mock.ExpectQuery("SELECT categories.id, categories.title, categories.slug, categories.created_at FROM categories").WillReturnRows(categoryRows)

	post, err := repo.FindById("1")

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "1", post.Id)
	assert.Equal(t, "Test Post", post.Title)
	assert.Len(t, post.Categories, 1)
}

func TestPostRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	mock.ExpectQuery("SELECT EXISTS").WithArgs("test-post").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery("INSERT INTO posts").WithArgs("Test Post", "test-post", "Content", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at"}).
			AddRow("1", "Test Post", "test-post", "Content", time.Now()))

	post := types.Post{
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Content",
		Author:  types.User{Id: "1"},
	}

	createdPost, err := repo.Create(post)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, "Test Post", createdPost.Title)
	assert.Equal(t, "test-post", createdPost.Slug)
}

func TestPostRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	// Mock finding the existing post
	mock.ExpectQuery("SELECT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email FROM posts").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email"}).
			AddRow("1", "Old Title", "old-slug", "Old Content", time.Now(), "1", "John", "Doe", "john@example.com"))

	// Mock fetching categories for the existing post
	mock.ExpectQuery("SELECT categories.id, categories.title, categories.slug, categories.created_at FROM categories").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}))

	// Mock checking if the new slug exists
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("new-slug").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Mock updating the post
	mock.ExpectQuery("UPDATE posts").
		WithArgs("New Title", "new-slug", "New Content", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at"}).
			AddRow("1", "New Title", "new-slug", "New Content", time.Now()))

	post := types.Post{
		Title:   "New Title",
		Slug:    "new-slug",
		Content: "New Content",
	}

	updatedPost, err := repo.Update("1", post)

	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, "New Title", updatedPost.Title)
	assert.Equal(t, "new-slug", updatedPost.Slug)
}
func TestPostRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	mock.ExpectExec("DELETE FROM posts").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete("1")

	assert.NoError(t, err)
}

func TestPostRepository_AssignCategoryToPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	mock.ExpectExec("INSERT INTO post_categories").WithArgs("1", "2").WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AssignCategoryToPost("1", "2")

	assert.NoError(t, err)
}

func TestPostRepository_UnassignCategoryFromPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	mock.ExpectExec("DELETE FROM post_categories").
		WithArgs("2", "1"). // Swap the order: categoryId first, then postId
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UnassignCategoryFromPost("1", "2")

	assert.NoError(t, err)
}

func TestPostRepository_GetCategoriesForPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "created_at"}).
		AddRow("1", "Category 1", "category-1", time.Now()).
		AddRow("2", "Category 2", "category-2", time.Now())

	mock.ExpectQuery("SELECT c.id, c.title, c.slug, c.created_at FROM categories c").WithArgs("1").WillReturnRows(rows)

	categories, err := repo.GetCategoriesForPost("1")

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Category 1", categories[0].Title)
	assert.Equal(t, "Category 2", categories[1].Title)
}

func TestPostRepository_UpdatePostCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM post_categories").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO post_categories").WithArgs("1", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO post_categories").WithArgs("1", "3").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdatePostCategories("1", []string{"2", "3"})

	assert.NoError(t, err)
}
func TestPostRepository_CreateWithUniqueSlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	// First attempt: "test-slug" already exists
	mock.ExpectQuery("SELECT EXISTS").WithArgs("test-slug").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	// Second attempt: "test-slug-1" is available
	mock.ExpectQuery("SELECT EXISTS").WithArgs("test-slug-1").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery("INSERT INTO posts").WithArgs("Test Post", "test-slug-1", "Content", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at"}).
			AddRow("1", "Test Post", "test-slug-1", "Content", time.Now()))

	post := types.Post{
		Title:   "Test Post",
		Slug:    "test-slug",
		Content: "Content",
		Author:  types.User{Id: "1"},
	}

	createdPost, err := repo.Create(post)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, "test-slug-1", createdPost.Slug)
}

func TestPostRepository_CreateWithMultipleSlugAttempts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	// First three attempts return true (slug exists)
	for i := 0; i < 3; i++ {
		slug := fmt.Sprintf("test-slug%s", func() string {
			if i == 0 {
				return ""
			}
			return fmt.Sprintf("-%d", i)
		}())
		mock.ExpectQuery("SELECT EXISTS").WithArgs(slug).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	}

	// Fourth attempt succeeds
	mock.ExpectQuery("SELECT EXISTS").WithArgs("test-slug-3").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery("INSERT INTO posts").WithArgs("Test Post", "test-slug-3", "Content", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at"}).
			AddRow("1", "Test Post", "test-slug-3", "Content", time.Now()))

	post := types.Post{
		Title:   "Test Post",
		Slug:    "test-slug",
		Content: "Content",
		Author:  types.User{Id: "1"},
	}

	createdPost, err := repo.Create(post)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, "test-slug-3", createdPost.Slug)
}

func TestPostRepository_ErrorHandling(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostRepository(db)

	t.Run("FindAll Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email, categories.id, categories.title, categories.slug, categories.created_at FROM posts").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.FindAll()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing FindAll query")
	})

	t.Run("FindAllPaginated Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").WillReturnError(sql.ErrConnDone)

		_, _, err := repo.FindAllPaginated(1, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing count query")
	})

	t.Run("FindBySlug Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email FROM posts").
			WithArgs("non-existent-slug").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.FindBySlug("non-existent-slug")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error scanning row in FindBySlug")
	})

	t.Run("FindById Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email FROM posts").
			WithArgs("non-existent-id").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.FindById("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error scanning row in FindById")
	})

	t.Run("Create Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").WithArgs("test-slug").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectQuery("INSERT INTO posts").WillReturnError(errors.New("duplicate key value violates unique constraint"))

		post := types.Post{
			Title:   "Test Post",
			Slug:    "test-slug",
			Content: "Content",
			Author:  types.User{Id: "1"},
		}

		_, err := repo.Create(post)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing Create query")
	})

	t.Run("Update Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email FROM posts").
			WithArgs("1").
			WillReturnError(sql.ErrNoRows)

		post := types.Post{
			Title:   "Updated Post",
			Slug:    "updated-slug",
			Content: "Updated Content",
		}

		_, err := repo.Update("1", post)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error finding post to update")
	})

	t.Run("Delete Error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM posts").WithArgs("non-existent-id").WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no post found with id non-existent-id to delete")
	})

	t.Run("AssignCategoryToPost Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO post_categories").WithArgs("1", "2").WillReturnError(errors.New("foreign key constraint violation"))

		err := repo.AssignCategoryToPost("1", "2")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing AssignCategoryToPost query")
	})

	t.Run("UnassignCategoryFromPost Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repository.NewPostRepository(db)

		mock.ExpectExec("DELETE FROM post_categories").
			WithArgs("2", "1"). // Swap the order: categoryId first, then postId
			WillReturnError(sql.ErrConnDone)

		err = repo.UnassignCategoryFromPost("1", "2")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing UnassignCategoryFromPost query")
	})

	t.Run("GetCategoriesForPost Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT c.id, c.title, c.slug, c.created_at FROM categories c").
			WithArgs("1").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetCategoriesForPost("1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing GetCategoriesForPost query")
	})

	t.Run("UpdatePostCategories Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := repository.NewPostRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM post_categories").
			WithArgs("1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO post_categories").
			WithArgs("1", "2").
			WillReturnError(errors.New("foreign key constraint violation"))
		mock.ExpectRollback()

		err = repo.UpdatePostCategories("1", []string{"2", "3"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error adding new category")
	})
}
