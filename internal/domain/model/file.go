package model

import (
	"time"
)

type File struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mimeType"`
	UploaderID uint      `json:"uploaderId"`
	UploadedAt time.Time `json:"uploadedAt"`
}

func NewFile(id, name, path string, size int64, mimeType string, uploaderID uint) *File {
	return &File{
		ID:         id,
		Name:       name,
		Path:       path,
		Size:       size,
		MimeType:   mimeType,
		UploaderID: uploaderID,
		UploadedAt: time.Now(),
	}
}
