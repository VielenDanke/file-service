package repository

import (
	"context"

	"github.com/vielendanke/file-service/internal/app/fileservice/model"
)

// FileRepository ...
type FileRepository interface {
	SaveFile(ctx context.Context, f model.FileModel, metadata string) error
	FindFileMetadataByID(ctx context.Context, id string) (map[string]string, error)
	FindFileNameByID(ctx context.Context, id string) (string, error)
	UpdateFileMetadataByID(ctx context.Context, metadata, id string) error
}
