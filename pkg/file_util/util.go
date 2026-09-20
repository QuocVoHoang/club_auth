package file_util

import (
	"mime"
	"path/filepath"

	"github.com/google/uuid"
)

// GetContentType detects content type from filename
func GetContentType(filename string) string {
	// Get content type based on file extension
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)

	// If mime type not found, use default
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType
}

// GenerateObjectName generates a unique object name using UUID + original file extension
func GenerateObjectName(originalFilename string) string {
	// Generate UUID v4
	uuidStr := uuid.New().String()

	// Keep the original file extension
	fileExtension := filepath.Ext(originalFilename)
	return uuidStr + fileExtension
}
