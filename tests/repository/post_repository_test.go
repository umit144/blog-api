package repository_test

import (
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

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email"}).
		AddRow("1", "Test Post", "test-post", "Content", time.Now(), "1", "John", "Doe", "john@example.com").
		AddRow("2", "Another Post", "another-post", "More Content", time.Now(), "2", "Jane", "Doe", "jane@example.com")

	mock.ExpectQuery("SELECT (.+) FROM posts").WillReturnRows(rows)

	posts, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Test Post", posts[0].Title)
	assert.Equal(t, "Another Post", posts[1].Title)
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

	mock.ExpectQuery("SELECT (.+) FROM posts").WithArgs("test-post").WillReturnRows(rows)

	post, err := repo.FindBySlug("test-post")

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "test-post", post.Slug)
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

	mock.ExpectQuery("SELECT (.+) FROM posts").WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "created_at", "user_id", "name", "lastname", "email"}).
			AddRow("1", "Old Title", "old-slug", "Old Content", time.Now(), "1", "John", "Doe", "john@example.com"))

	mock.ExpectQuery("SELECT EXISTS").WithArgs("new-slug").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery("UPDATE posts").WithArgs("New Title", "new-slug", "New Content", "1").
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
