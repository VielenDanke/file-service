package repository

import (
	"context"

	"github.com/vielendanke/file-service/internal/app/fileservice/model"
)

// FileRepository ...
type FileRepository interface {
	SaveFileMetadata(ctx context.Context, f model.FileModel, metadata string) error
	FindFileMetadataByID(ctx context.Context, id string) (map[string]string, error)
	FindFileNameByID(ctx context.Context, id string) (string, error)
	UpdateFileMetadataByID(ctx context.Context, metadata, id string) error
	CheckIfExists(ctx context.Context, f model.FileModel) error
	DeleteMetadataByID(ctx context.Context, id string) error
}
