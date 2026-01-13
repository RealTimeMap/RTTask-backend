package dto

import (
	"fmt"
	"rttask/internal/domain/model"
)

type FileResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
	Url      string `json:"url"`
}

func NewFileResponse(file *model.File) *FileResponse {
	if file == nil {
		return nil
	}
	return &FileResponse{
		ID:       file.ID,
		Name:     file.Name,
		Size:     file.Size,
		MimeType: file.MimeType,
		Url:      fmt.Sprintf("%s://%s/store/%s", "https", "realtimemap.ru/rttask", file.Path),
	}
}

// NewMultiFileResponse создает массив FileResponse из массива указателей на File.
// Пропускает nil элементы, возвращает пустой массив если входной массив пустой или nil
func NewMultiFileResponse(files []*model.File) []FileResponse {
	if files == nil || len(files) == 0 {
		return []FileResponse{}
	}

	response := make([]FileResponse, 0, len(files))
	for _, file := range files {
		if fileResp := NewFileResponse(file); fileResp != nil {
			response = append(response, *fileResp)
		}
	}
	return response
}
