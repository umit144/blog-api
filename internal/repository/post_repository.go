package repository

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go-blog/internal/types"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (repo PostRepository) FindAll() ([]types.Post, error) {
	var posts []types.Post

	sql, args, err := sq.Select("posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email").
		From("posts").
		Join("users ON posts.user_id = users.id").
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
		var post types.Post
		var user types.User
		err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.CreatedAt,
			&user.Id,
			&user.Name,
			&user.Lastname,
			&user.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row in FindAll: %v", err)
		}
		post.Author = user
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows in FindAll: %v", err)
	}

	return posts, nil
}

func (repo PostRepository) FindBySlug(slug string) (*types.Post, error) {
	query := sq.Select("posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email").
		From("posts").
		Join("users ON posts.user_id = users.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"posts.slug": slug})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindBySlug: %v", err)
	}

	var post types.Post
	var user types.User
	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&post.Id,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.CreatedAt,
		&user.Id,
		&user.Name,
		&user.Lastname,
		&user.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing FindBySlug query: %v", err)
	}

	post.Author = user
	return &post, nil
}

func (repo PostRepository) FindById(id string) (*types.Post, error) {
	query := sq.Select("posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email").
		From("posts").
		Join("users ON posts.user_id = users.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"posts.id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for FindById: %v", err)
	}

	var post types.Post
	var user types.User
	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&post.Id,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.CreatedAt,
		&user.Id,
		&user.Name,
		&user.Lastname,
		&user.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing FindById query: %v", err)
	}

	post.Author = user
	return &post, nil
}

// Create, Update, Delete ve generateUniqueSlug fonksiyonları aynı kalacak

func (repo PostRepository) Create(post types.Post) (*types.Post, error) {
	slug, err := repo.generateUniqueSlug(post.Slug, "")
	if err != nil {
		return nil, fmt.Errorf("error generating unique slug: %v", err)
	}

	insertQuery := sq.Insert("posts").
		Columns("title", "slug", "content", "user_id").
		Values(post.Title, slug, post.Content, post.Author.Id).
		Suffix("RETURNING id, title, slug, content, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for Create: %v", err)
	}

	var createdPost types.Post

	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&createdPost.Id,
		&createdPost.Title,
		&createdPost.Slug,
		&createdPost.Content,
		&createdPost.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing Create query: %v", err)
	}

	createdPost.Author = post.Author
	return &createdPost, nil
}

func (repo PostRepository) Update(id string, post types.Post) (*types.Post, error) {
	existingPost, err := repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("error finding post to update: %v", err)
	}

	slug, err := repo.generateUniqueSlug(post.Slug, id)
	if err != nil {
		return nil, fmt.Errorf("error generating unique slug for update: %v", err)
	}

	updateQuery := sq.Update("posts").
		Set("title", post.Title).
		Set("slug", slug).
		Set("content", post.Content).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, title, slug, content, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := updateQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for Update: %v", err)
	}

	var updatedPost types.Post
	err = repo.db.QueryRowContext(context.Background(), sql, args...).Scan(
		&updatedPost.Id,
		&updatedPost.Title,
		&updatedPost.Slug,
		&updatedPost.Content,
		&updatedPost.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error executing Update query: %v", err)
	}

	updatedPost.Author = existingPost.Author

	return &updatedPost, nil
}

func (repo PostRepository) Delete(id string) error {
	deleteQuery := sq.Delete("posts").
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
		return fmt.Errorf("no post found with id %s to delete", id)
	}

	return nil
}

func (repo PostRepository) generateUniqueSlug(baseSlug string, excludeId string) (string, error) {
	slug := baseSlug
	counter := 1

	for {
		query := sq.Select("EXISTS(SELECT 1 FROM posts WHERE slug = ? AND id != ?)").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := query.ToSql()
		if err != nil {
			return "", fmt.Errorf("error creating SQL for slug check: %v", err)
		}

		var exists bool
		err = repo.db.QueryRowContext(context.Background(), sql, append(args, slug, excludeId)...).Scan(&exists)
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
