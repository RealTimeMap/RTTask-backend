package file

import (
	"context"
	"fmt"
	"path/filepath"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/infrastructure/storage"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var mimeMap = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".txt":  "text/plain",
	".csv":  "text/csv",
	".zip":  "application/zip",
}

type FileService struct {
	storage storage.FileStorage
	logger  *zap.Logger
}

func NewFileService(storage storage.FileStorage, logger *zap.Logger) *FileService {
	return &FileService{
		storage: storage,
		logger:  logger,
	}
}

func (s *FileService) UploadFile(ctx context.Context, input FileInput, profile ValidationProfile) (*model.File, error) {
	if err := s.validate(input, profile); err != nil {
		return nil, err
	}
	fileID := uuid.New().String()
	filePath := s.generateFilePath(input.EntityType, fileID, input.FileHeader.Filename)

	if err := s.storage.Save(ctx, input.File, filePath); err != nil {
		return nil, err
	}
	newFile := model.NewFile(
		fileID,
		input.FileHeader.Filename,
		filePath,
		input.FileHeader.Size,
		input.FileHeader.Header.Get("Content-Type"),
		input.UploaderID,
	)
	return newFile, nil
}

func (s *FileService) generateFilePath(entityType string, fileID string, fileName string) string {
	now := time.Now()
	ext := filepath.Ext(fileName)

	// Формат: {entity_type}/{year}/{month}/{uuid}{ext}
	return fmt.Sprintf("%s/%d/%02d/%s%s", entityType, now.Year(), now.Month(), fileID, ext)
}

func (s *FileService) validate(input FileInput, profile ValidationProfile) error {
	if input.FileHeader.Size > profile.MaxFileSize {
		return domainerrors.NewValidationError("File size too big")
	}
	if input.FileHeader.Size == 0 {
		return domainerrors.NewValidationError("File is empty")
	}
	mimiType := input.FileHeader.Header.Get("Content-Type")
	if mimiType == "" {
		mimiType = s.mimeTypeByExtension(input.FileHeader.Filename)
	}
	if !s.isAllowedMime(mimiType, profile.AllowedMimes) {
		return domainerrors.NewValidationError("Mime type not allowed")
	}
	return nil
}

func (s *FileService) mimeTypeByExtension(fileName string) string {
	ext := filepath.Ext(fileName)

	if mime, ok := mimeMap[ext]; ok {
		return mime
	}
	return ""
}

func (s *FileService) isAllowedMime(mime string, allowedMimes []string) bool {
	mime = strings.ToLower(mime)
	for _, allowed := range allowedMimes {
		if strings.ToLower(allowed) == mime {
			return true
		}
	}
	return false
}
