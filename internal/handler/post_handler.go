package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-blog/internal/database"
	"go-blog/internal/repository"
	"go-blog/internal/types"
	"regexp"
	"strings"
)

type PostHandler struct {
	postRepository repository.PostRepository
}

func NewPostHandler(db database.Service) *PostHandler {
	return &PostHandler{
		postRepository: *repository.NewPostRepository(db.GetInstance()),
	}
}

func (h *PostHandler) GetPostHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if slug != "" {
		post, err := h.postRepository.FindBySlug(slug)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Post not found",
				"message": fmt.Sprintf("No post found with slug: %s", slug),
			})
		}
		return c.JSON(post)
	}

	posts, err := h.postRepository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve posts",
			"message": fmt.Sprintf("Error occurred while fetching posts: %v", err),
		})
	}

	return c.JSON(posts)
}

func (h *PostHandler) CreatePostHandler(c *fiber.Ctx) error {
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

	user, ok := c.Locals("user").(*types.User)
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

func (h *PostHandler) UpdatePostHandler(c *fiber.Ctx) error {
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

	user, ok := c.Locals("user").(*types.User)
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

func (h *PostHandler) DeletePostHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	user, ok := c.Locals("user").(*types.User)
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
