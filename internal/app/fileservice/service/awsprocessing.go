package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	s3store "github.com/unistack-org/micro-store-s3/v3"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/store"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
	"github.com/vielendanke/file-service/internal/app/fileservice/repository"
)

var keyRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// AWSProcessingService ...
type AWSProcessingService struct {
	fileRepository repository.FileRepository
	codec          codec.Codec
	cleanStore     store.Store
	dirtyStore     store.Store
}

// NewAWSProcessingService ...
func NewAWSProcessingService(codec codec.Codec, fileRepository repository.FileRepository, cleanStore store.Store, dirtyStore store.Store) FileProcessingService {
	return &AWSProcessingService{
		fileRepository: fileRepository,
		codec:          codec,
		cleanStore:     cleanStore,
		dirtyStore:     dirtyStore,
	}
}

// StoreFile ...
func (aps *AWSProcessingService) StoreFile(ctx context.Context, f model.FileModel) error {
	awsFile := f.(*model.AWSModel)
	fileID := uuid.New().String()
	fileID = strings.ReplaceAll(fileID, "-", "")
	awsFile.FileID = fileID
	if err := aps.dirtyStore.Write(
		ctx,
		fileID,
		awsFile.File,
		s3store.WriteBucket("micro-store-s3"),
		s3store.ContentType("application/octet-stream"),
	); err != nil {
		return fmt.Errorf("Error writing file to s3, %v", err)
	}
	if err := aps.dirtyStore.Exists(ctx, fileID, s3store.ExistsBucket("micro-store-s3")); err != nil {
		return fmt.Errorf("File not exists in store, %v", err)
	}
	return nil
}

// SaveFileData ...
func (aps *AWSProcessingService) SaveFileData(ctx context.Context, f model.FileModel) error {
	awsFile := f.(*model.AWSModel)
	PrepareMetadata(awsFile.Metadata, []string{"type", "class", "number"})
	jsonMetadata, err := aps.codec.Marshal(awsFile.GetMetadata())
	if err != nil {
		return fmt.Errorf("Error while marshalling metadata, %v", err)
	}
	if len(awsFile.GetFileID()) == 0 {
		return fmt.Errorf("ID is not present")
	}
	if err := aps.fileRepository.SaveFile(ctx, awsFile, string(jsonMetadata)); err != nil {
		return err
	}
	return nil
}

// GetFileMetadata ...
func (aps *AWSProcessingService) GetFileMetadata(ctx context.Context, id string) (map[string]interface{}, error) {
	properties, err := aps.fileRepository.FindFileMetadataByID(ctx, id)
	if err != nil {
		return nil, err
	}
	metadata := properties["metadata"]
	jsonMap := make(map[string]interface{})
	if err := aps.codec.Unmarshal([]byte(metadata), &jsonMap); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSONB, %v", err)
	}
	jsonMap["type"] = properties["type"]
	jsonMap["class"] = properties["class"]
	jsonMap["number"] = properties["number"]
	return jsonMap, nil
}

// DownloadFile ...
func (aps *AWSProcessingService) DownloadFile(ctx context.Context, id string) ([]byte, string, error) {
	wg := sync.WaitGroup{}
	errCh := make(chan error, 2)
	defer close(errCh)
	file := []byte{}
	filename := ""
	wg.Add(2)
	go func(file *[]byte) {
		defer wg.Done()
		if err := aps.cleanStore.Read(ctx, id, file, s3store.ReadBucket("micro-store-s3")); err != nil {
			errCh <- fmt.Errorf("Error download file from s3, %v", err)
		}
	}(&file)
	go func(filename *string) {
		defer wg.Done()
		fn, err := aps.fileRepository.FindFileNameByID(ctx, id)
		if err != nil {
			errCh <- err
		}
		*filename = fn
	}(&filename)
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, "", err
	default:
		return file, filename, nil
	}
}

// UpdateFileMetadata ...
func (aps *AWSProcessingService) UpdateFileMetadata(ctx context.Context, metadata map[string]interface{}, id string) error {
	jsonMetadata, err := aps.codec.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("Error marshalling metadata: %v", err)
	}
	if err := aps.fileRepository.UpdateFileMetadataByID(ctx, string(jsonMetadata), id); err != nil {
		return err
	}
	return nil
}
