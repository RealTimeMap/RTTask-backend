package file

import (
	"fmt"
	"mime/multipart"
)

type FileInput struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
	EntityType string
	UploaderID uint
}

// NewFileInput создает FileInput из multipart.FileHeader
// Открывает файл и возвращает готовую структуру для FileService
func NewFileInput(
	fileHeader *multipart.FileHeader,
	entityType string,
	uploaderID uint,
) (FileInput, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return FileInput{}, fmt.Errorf("failed to open file: %w", err)
	}

	return FileInput{
		File:       file,
		FileHeader: fileHeader,
		EntityType: entityType,
		UploaderID: uploaderID,
	}, nil
}
