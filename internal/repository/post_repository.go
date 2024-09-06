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

	sql, _, _ := sq.Select("posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email").
		From("posts").
		Join("users ON posts.user_id = users.id").
		ToSql()

	rows, err := repo.db.QueryContext(context.Background(), sql)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		post.Author = user
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("error creating SQL: %v", err)
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
		return nil, err
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
		return nil, fmt.Errorf("error creating SQL: %v", err)
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
		return nil, err
	}

	post.Author = user
	return &post, nil
}

func (repo PostRepository) Create(post types.Post) (*types.Post, error) {
	baseSlug := post.Slug
	slug := baseSlug
	counter := 1

	for {
		query := sq.Select("EXISTS(SELECT 1 FROM posts WHERE slug = ?)").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := query.ToSql()
		if err != nil {
			return nil, err
		}

		var exists bool
		err = repo.db.QueryRowContext(context.Background(), sql, append(args, slug)...).Scan(&exists)
		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	insertQuery := sq.Insert("posts").
		Columns("title", "slug", "content", "user_id").
		Values(post.Title, slug, post.Content, post.Author.Id).
		Suffix("RETURNING id, title, slug, content, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &createdPost, nil
}

func (repo PostRepository) Update(id string, post types.Post) (*types.Post, error) {
	existingPost, err := repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("post not found: %v", err)
	}

	slug := post.Slug
	if slug != existingPost.Slug {
		baseSlug := slug
		counter := 1
		for {
			query := sq.Select("EXISTS(SELECT 1 FROM posts WHERE slug = ? AND id != ?)").
				PlaceholderFormat(sq.Dollar)

			sql, args, err := query.ToSql()
			if err != nil {
				return nil, err
			}

			var exists bool
			err = repo.db.QueryRowContext(context.Background(), sql, append(args, slug, id)...).Scan(&exists)
			if err != nil {
				return nil, err
			}

			if !exists {
				break
			}

			slug = fmt.Sprintf("%s-%d", baseSlug, counter)
			counter++
		}
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
		return nil, fmt.Errorf("error creating SQL: %v", err)
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
		return nil, fmt.Errorf("update error: %v", err)
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
		return fmt.Errorf("error creating SQL: %v", err)
	}

	result, err := repo.db.ExecContext(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("delete error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found to delete")
	}

	return nil
}
