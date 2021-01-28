package service

import (
	"context"

	"github.com/vielendanke/file-service/internal/app/fileservice/model"
)

// FileProcessingService ...
type FileProcessingService interface {
	StoreFile(ctx context.Context, f model.FileModel) error
	SaveFileData(ctx context.Context, f model.FileModel) error
	GetFileMetadata(ctx context.Context, id string) (map[string]interface{}, error)
	DownloadFile(ctx context.Context, id string) ([]byte, string, error)
	UpdateFileMetadata(ctx context.Context, metadata map[string]interface{}, id string) error
}
