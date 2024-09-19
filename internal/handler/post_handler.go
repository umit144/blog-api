package handler

import (
	"fmt"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PostHandler interface {
	GetPostHandler(c *fiber.Ctx) error
	CreatePostHandler(c *fiber.Ctx) error
	UpdatePostHandler(c *fiber.Ctx) error
	DeletePostHandler(c *fiber.Ctx) error
	AssignCategoryToPostHandler(c *fiber.Ctx) error
	UnassignCategoryFromPostHandler(c *fiber.Ctx) error
	GetCategoriesForPostHandler(c *fiber.Ctx) error
	UpdatePostCategoriesHandler(c *fiber.Ctx) error
}

type postHandler struct {
	postRepository repository.PostRepository
}

func NewPostHandler(postRepository repository.PostRepository) PostHandler {
	return &postHandler{postRepository}
}

func (h *postHandler) GetPostHandler(c *fiber.Ctx) error {
	slugOrId := c.Params("slugOrId")

	if slugOrId != "" {
		var post *types.Post

		_, err := uuid.Parse(slugOrId)
		if err == nil {
			post, err = h.postRepository.FindById(slugOrId)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Post not found",
					"message": fmt.Sprintf("No post found with id: %s", slugOrId),
				})
			}
		} else {
			post, err = h.postRepository.FindBySlug(slugOrId)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Post not found",
					"message": fmt.Sprintf("No post found with slug: %s", slugOrId),
				})
			}
		}

		return c.JSON(post)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	posts, totalCount, err := h.postRepository.FindAllPaginated(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve posts",
			"message": fmt.Sprintf("Error occurred while fetching posts: %v", err),
		})
	}

	totalPages := (totalCount + limit - 1) / limit

	return c.JSON(fiber.Map{
		"meta": fiber.Map{
			"page":       page,
			"limit":      limit,
			"totalCount": totalCount,
			"totalPages": totalPages,
		},
		"data": posts,
	})
}

func (h *postHandler) CreatePostHandler(c *fiber.Ctx) error {
	var post types.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing post data: %v", err),
		})
	}

	if err := post.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	user, ok := c.Locals("user").(types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	post.Author.Id = user.Id

	// Generate slug
	slug := strings.ToLower(post.Title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	post.Slug = strings.Trim(slug, "-")

	createdPost, err := h.postRepository.Create(post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create post",
			"message": fmt.Sprintf("Error occurred while creating post: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPost)
}

func (h *postHandler) UpdatePostHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var post types.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing post data: %v", err),
		})
	}

	if err := post.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"fails": err,
		})
	}

	user, ok := c.Locals("user").(types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	existingPost, err := h.postRepository.FindById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Post not found",
			"message": fmt.Sprintf("No post found with ID: %s", id),
		})
	}

	if existingPost.Author.Id != user.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this post",
		})
	}

	// Generate slug
	slug := strings.ToLower(post.Title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	post.Slug = strings.Trim(slug, "-")

	updatedPost, err := h.postRepository.Update(id, post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update post",
			"message": fmt.Sprintf("Error occurred while updating post: %v", err),
		})
	}

	return c.JSON(updatedPost)
}

func (h *postHandler) DeletePostHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, ok := c.Locals("user").(types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	existingPost, err := h.postRepository.FindById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Post not found",
			"message": fmt.Sprintf("No post found with ID: %s", id),
		})
	}

	if existingPost.Author.Id != user.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this post",
		})
	}

	if err := h.postRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete post",
			"message": fmt.Sprintf("Error occurred while deleting post: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *postHandler) AssignCategoryToPostHandler(c *fiber.Ctx) error {
	postId := c.Params("postId")
	categoryId := c.Params("categoryId")

	err := h.postRepository.AssignCategoryToPost(postId, categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to assign category to post",
			"message": fmt.Sprintf("Error occurred while assigning category: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *postHandler) UnassignCategoryFromPostHandler(c *fiber.Ctx) error {
	postId := c.Params("postId")
	categoryId := c.Params("categoryId")

	err := h.postRepository.UnassignCategoryFromPost(postId, categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to unassign category from post",
			"message": fmt.Sprintf("Error occurred while unassigning category: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *postHandler) GetCategoriesForPostHandler(c *fiber.Ctx) error {
	postId := c.Params("postId")

	categories, err := h.postRepository.GetCategoriesForPost(postId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get categories for post",
			"message": fmt.Sprintf("Error occurred while fetching categories: %v", err),
		})
	}

	return c.JSON(categories)
}

func (h *postHandler) UpdatePostCategoriesHandler(c *fiber.Ctx) error {
	postId := c.Params("postId")

	var request struct {
		CategoryIds []string `json:"categoryIds"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"message": fmt.Sprintf("Error parsing request body: %v", err),
		})
	}

	err := h.postRepository.UpdatePostCategories(postId, request.CategoryIds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update post categories",
			"message": fmt.Sprintf("Error occurred while updating categories: %v", err),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
