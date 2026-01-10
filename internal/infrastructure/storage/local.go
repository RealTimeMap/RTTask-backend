package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	domainerrors "rttask/internal/domain/errors"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) FileStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}
func (s *LocalStorage) Save(ctx context.Context, file io.Reader, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return domainerrors.NewValidationError("make dir error")
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return domainerrors.NewValidationError("create file error")
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		os.Remove(fullPath)
		return domainerrors.NewValidationError("copy file error")
	}
	return nil
}
