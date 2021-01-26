package service

import "github.com/vielendanke/file-service/internal/app/fileservice/model"

// FileProcessingService ...
type FileProcessingService interface {
	StoreFile(f model.FileModel) error
	SaveFileData(f model.FileModel) (string, error)
	GetFileMetadata(id string) (string, error)
}
