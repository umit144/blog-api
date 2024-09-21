package handler

import (
	"go-blog/internal/service"
	"go-blog/internal/types"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type fileHandler struct {
	fileService service.FileService
}

type FileHandler interface {
	UploadFileHandler(c *fiber.Ctx) error
	DeleteFileHandler(c *fiber.Ctx) error
}

func NewFileHandler(fileService service.FileService) FileHandler {
	return &fileHandler{fileService}
}

func (h *fileHandler) UploadFileHandler(c *fiber.Ctx) error {
	var user = c.Locals("user").(types.User)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File upload failed",
		})
	}

	filename := filepath.Base(file.Filename)
	uniqueFilename, err := h.fileService.GenerateUniqueFilename(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate unique filename",
		})
	}

	err = h.fileService.SaveFile(file, uniqueFilename, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "File uploaded successfully",
		"filename": uniqueFilename,
	})
}

func (h *fileHandler) DeleteFileHandler(c *fiber.Ctx) error {
	var user = c.Locals("user").(types.User)
	filename := c.Params("filename")

	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Filename is required",
		})
	}

	err := h.fileService.DeleteFile(filename, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}
