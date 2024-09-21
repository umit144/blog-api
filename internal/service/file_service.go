package service

import (
	"fmt"
	"go-blog/internal/types"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileService interface {
	GenerateUniqueFilename(filename string) (string, error)
	SaveFile(file *multipart.FileHeader, filename string, user types.User) error
	DeleteFile(filename string, user types.User) error
}

type fileService struct {
	uploadDir string
}

func NewFileService() FileService {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current working directory: %v", err))
	}

	uploadDir := filepath.Join(rootDir, "uploads")
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create uploads directory: %v", err))
	}

	return &fileService{uploadDir: uploadDir}
}

func (s *fileService) GenerateUniqueFilename(filename string) (string, error) {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	timestamp := time.Now().Format("20060102150405")
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s_%s%s", name, timestamp, uniqueID, ext), nil
}

func (s *fileService) SaveFile(file *multipart.FileHeader, filename string, user types.User) error {
	userDir := filepath.Join(s.uploadDir, user.Id)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return fmt.Errorf("failed to create user directory: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(userDir, filename))
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

func (s *fileService) DeleteFile(filename string, user types.User) error {
	userDir := filepath.Join(s.uploadDir, user.Id)
	filePath := filepath.Join(userDir, filename)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %w", err)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
