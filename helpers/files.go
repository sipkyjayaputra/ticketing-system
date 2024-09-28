package helpers

import (
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveUploadedFile(fileHeader *multipart.FileHeader) (string, error) {
	// Define the directory to save files
	uploadDir := "./uploads/"
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create the file path
	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	// Save the file to disk
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
