package handler

import (
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
				"error": "Post not found",
			})
		}
		return c.JSON(post)
	}

	posts, err := h.postRepository.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(posts)
}

func (h *PostHandler) CreatePostHandler(c *fiber.Ctx) error {
	var post types.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
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

	slug := strings.ToLower(post.Title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	post.Slug = strings.Trim(slug, "-")

	createdPost, err := h.postRepository.Create(post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPost)
}

func (h *PostHandler) UpdatePostHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var post types.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
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
			"error": "Post not found",
		})
	}

	if existingPost.Author.Id != user.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this post",
		})
	}

	slug := strings.ToLower(post.Title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	post.Slug = strings.Trim(slug, "-")

	updatedPost, err := h.postRepository.Update(id, post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update post",
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
			"error": "Post not found",
		})
	}

	if existingPost.Author.Id != user.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this post",
		})
	}

	if err := h.postRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete post",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
