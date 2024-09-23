package service_test

import (
	"bytes"
	"go-blog/internal/service"
	"go-blog/internal/types"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()

	err := os.RemoveAll("uploads")
	if err != nil {
		panic("failed to clean up uploads directory: " + err.Error())
	}

	os.Exit(code)
}

func TestNewFileService(t *testing.T) {
	fs := service.NewFileService()
	assert.NotNil(t, fs, "NewFileService should return a non-nil FileService")
}

func TestGenerateUniqueFilename(t *testing.T) {
	fs := service.NewFileService()
	filename := "test.txt"

	uniqueFilename, err := fs.GenerateUniqueFilename(filename)
	require.NoError(t, err)
	assert.Contains(t, uniqueFilename, "test_")
	assert.Contains(t, uniqueFilename, ".txt")
	assert.Len(t, strings.Split(uniqueFilename, "_"), 3)
}

func TestSaveFile(t *testing.T) {
	fs := service.NewFileService()
	user := types.User{Id: "testuser"}

	fileContent := []byte("test content")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)

	_, err = part.Write(fileContent)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(32 << 20)
	require.NoError(t, err)

	file, header, err := req.FormFile("file")
	require.NoError(t, err)
	defer file.Close()

	err = fs.SaveFile(header, "test.txt", user)
	require.NoError(t, err)

	savedFilePath := filepath.Join("uploads", user.Id, "test.txt")
	assert.FileExists(t, savedFilePath)

	savedContent, err := os.ReadFile(savedFilePath)
	require.NoError(t, err)
	assert.Equal(t, fileContent, savedContent)
}

func TestDeleteFile(t *testing.T) {
	fs := service.NewFileService()
	user := types.User{Id: "testuser"}

	testDir := filepath.Join("uploads", user.Id)
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	testFilePath := filepath.Join(testDir, "test.txt")
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	err = fs.DeleteFile("test.txt", user)
	require.NoError(t, err)

	_, err = os.Stat(testFilePath)
	assert.True(t, os.IsNotExist(err))

	err = fs.DeleteFile("nonexistent.txt", user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}
