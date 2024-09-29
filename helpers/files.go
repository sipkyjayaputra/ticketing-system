package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveUploadedFile uploads the form file to the specified destination.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a destination file where the uploaded file will be saved
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the contents of the uploaded file to the destination file
	_, err = io.Copy(out, src)
	return err
}
