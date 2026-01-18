package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	Save(ctx context.Context, file io.Reader, path string) error
	Delete(ctx context.Context, path string) error
}
