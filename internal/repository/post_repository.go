package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-blog/internal/types"

	sq "github.com/Masterminds/squirrel"
)

type PostRepository interface {
	FindAll() ([]types.Post, error)
	FindAllPaginated(page, limit int) ([]types.Post, int, error)
	FindBySlug(slug string) (*types.Post, error)
	FindById(id string) (*types.Post, error)
	Create(post types.Post) (*types.Post, error)
	Update(id string, post types.Post) (*types.Post, error)
	Delete(id string) error
	AssignCategoryToPost(postId string, categoryId string) error
	UnassignCategoryFromPost(postId string, categoryId string) error
	GetCategoriesForPost(postId string) ([]types.Category, error)
	UpdatePostCategories(postId string, categoryIds []string) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (repo postRepository) FindAll() ([]types.Post, error) {
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

func (repo postRepository) FindAllPaginated(page, limit int) ([]types.Post, int, error) {
	var posts []types.Post
	var totalCount int

	// First, get the total count of posts
	countSql, countArgs, err := sq.Select("COUNT(*)").
		From("posts").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("error creating SQL for count query: %v", err)
	}

	err = repo.db.QueryRowContext(context.Background(), countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing count query: %v", err)
	}

	// Now, get the paginated posts
	offset := (page - 1) * limit

	sql, args, err := sq.Select("posts.id, posts.title, posts.slug, posts.content, posts.created_at, users.id, users.name, users.lastname, users.email").
		From("posts").
		Join("users ON posts.user_id = users.id").
		OrderBy("posts.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("error creating SQL for FindAllPaginated: %v", err)
	}

	rows, err := repo.db.QueryContext(context.Background(), sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing FindAllPaginated query: %v", err)
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
			return nil, 0, fmt.Errorf("error scanning row in FindAllPaginated: %v", err)
		}
		post.Author = user
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error after iterating rows in FindAllPaginated: %v", err)
	}

	return posts, totalCount, nil
}

func (repo postRepository) FindBySlug(slug string) (*types.Post, error) {
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

func (repo postRepository) FindById(id string) (*types.Post, error) {
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

func (repo postRepository) Create(post types.Post) (*types.Post, error) {
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

func (repo postRepository) Update(id string, post types.Post) (*types.Post, error) {
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

func (repo postRepository) Delete(id string) error {
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

func (repo postRepository) generateUniqueSlug(baseSlug string, excludeId string) (string, error) {
	slug := baseSlug
	counter := 1

	for {
		query := sq.Select("EXISTS(SELECT 1 FROM posts WHERE slug = ?)").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := query.ToSql()
		if err != nil {
			return "", fmt.Errorf("error creating SQL for slug check: %v", err)
		}

		var exists bool
		err = repo.db.QueryRowContext(context.Background(), sql, append(args, slug)...).Scan(&exists)
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

func (repo postRepository) AssignCategoryToPost(postId string, categoryId string) error {
	query := sq.Insert("post_categories").
		Columns("post_id", "category_id").
		Values(postId, categoryId).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error creating SQL for AssignCategoryToPost: %v", err)
	}

	_, err = repo.db.ExecContext(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("error executing AssignCategoryToPost query: %v", err)
	}

	return nil
}

func (repo postRepository) UnassignCategoryFromPost(postId string, categoryId string) error {
	query := sq.Delete("post_categories").
		Where(sq.Eq{"post_id": postId, "category_id": categoryId}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error creating SQL for UnassignCategoryFromPost: %v", err)
	}

	_, err = repo.db.ExecContext(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("error executing UnassignCategoryFromPost query: %v", err)
	}

	return nil
}

func (repo postRepository) GetCategoriesForPost(postId string) ([]types.Category, error) {
	query := sq.Select("c.id", "c.title", "c.slug", "c.created_at").
		From("categories c").
		Join("post_categories pc ON c.id = pc.category_id").
		Where(sq.Eq{"pc.post_id": postId}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error creating SQL for GetCategoriesForPost: %v", err)
	}

	rows, err := repo.db.QueryContext(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing GetCategoriesForPost query: %v", err)
	}
	defer rows.Close()

	var categories []types.Category
	for rows.Next() {
		var category types.Category
		err := rows.Scan(&category.Id, &category.Title, &category.Slug, &category.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row in GetCategoriesForPost: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows in GetCategoriesForPost: %v", err)
	}

	return categories, nil
}

func (repo postRepository) UpdatePostCategories(postId string, categoryIds []string) error {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Remove all existing categories for the post
	_, err = tx.ExecContext(context.Background(), "DELETE FROM post_categories WHERE post_id = $1", postId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error removing existing categories: %v", err)
	}

	// Add new categories
	for _, categoryId := range categoryIds {
		_, err = tx.ExecContext(context.Background(), "INSERT INTO post_categories (post_id, category_id) VALUES ($1, $2)", postId, categoryId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error adding new category: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
