package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveUploadedFile(header *multipart.FileHeader, uploadDir string) (string, error) {
	// Create the upload dir if it dies not exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}
	}

	// Gnereate unique file name
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	fullpath := filepath.Join(uploadDir, filename)

	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create the destination file
	destination, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("failee to create file: %v", err)
	}

	defer destination.Close()

	// Copy the uploaded file to the destination
	if _, err := io.Copy(destination, src); err != nil {
		return "", fmt.Errorf("fialed to save file: %v", err)
	}

	// Return the path to the saved file
	return fullpath, nil
}
