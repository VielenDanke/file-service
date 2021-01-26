package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	s3store "github.com/unistack-org/micro-store-s3/v3"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/store"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
)

var keyRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// AWSProcessingService ...
type AWSProcessingService struct {
	db    *sqlx.DB
	codec codec.Codec
	store store.Store
}

// NewAWSProcessingService ...
func NewAWSProcessingService(codec codec.Codec, db *sqlx.DB, store store.Store) FileProcessingService {
	return &AWSProcessingService{
		db:    db,
		codec: codec,
		store: store,
	}
}

// StoreFile ...
func (aps *AWSProcessingService) StoreFile(ctx context.Context, f model.FileModel) error {
	awsFile := f.(*model.AWSModel)
	fileID := uuid.New().String()
	fileID = strings.ReplaceAll(fileID, "-", "")
	awsFile.FileID = fileID
	if err := aps.store.Write(
		ctx,
		fileID,
		awsFile.File,
		s3store.WriteBucket("micro-store-s3"),
		s3store.ContentType("application/octet-stream"),
	); err != nil {
		return fmt.Errorf("Error writing file to s3, %v", err)
	}
	if err := aps.store.Exists(ctx, fileID, s3store.ExistsBucket("micro-store-s3")); err != nil {
		return fmt.Errorf("File not exists in store, %v", err)
	}
	return nil
}

// SaveFileData ...
func (aps *AWSProcessingService) SaveFileData(ctx context.Context, f model.FileModel) error {
	awsFile := f.(*model.AWSModel)
	jsonMetadata, err := aps.codec.Marshal(awsFile.GetMetadata())
	if err != nil {
		return fmt.Errorf("Error while marshalling metadata, %v", err)
	}
	if len(awsFile.GetFileID()) == 0 {
		return fmt.Errorf("ID is not present")
	}
	tx := aps.db.MustBegin()
	res := tx.MustExecContext(ctx, "INSERT INTO FILES(ID, FILE_NAME, IIN, FILE_LINK, METADATA) VALUES($1, $2, $3, $4, $5)",
		awsFile.GetFileID(), awsFile.GetFileName(), awsFile.GetIIN(), awsFile.GetFileLink(), string(jsonMetadata),
	)
	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error during rows inserter, %v", err)
	}
	if num == 0 {
		return fmt.Errorf("No insertions found, %d", num)
	}
	tx.Commit()
	return nil
}

// GetFileMetadata ...
func (aps *AWSProcessingService) GetFileMetadata(ctx context.Context, id string) (string, error) {
	jsonMetadata := ""
	if err := aps.db.QueryRowContext(ctx, "SELECT METADATA FROM FILES WHERE ID=$1", id).Scan(&jsonMetadata); err != nil {
		return "", fmt.Errorf("Error while reading from DB, %v", err)
	}
	return jsonMetadata, nil
}

// DownloadFile ...
func (aps *AWSProcessingService) DownloadFile(ctx context.Context, id string) ([]byte, string, error) {
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	defer close(errCh)
	file := []byte{}
	filename := ""
	wg.Add(2)
	go func(file *[]byte) {
		defer wg.Done()
		if err := aps.store.Read(ctx, id, file, s3store.ReadBucket("micro-store-s3")); err != nil {
			errCh <- fmt.Errorf("Error download file from s3, %v", err)
		}
	}(&file)
	go func(filename *string) {
		defer wg.Done()
		if row := aps.db.QueryRowContext(ctx, "SELECT FILE_NAME FROM FILES WHERE ID=$1", id).Scan(filename); row != nil {
			errCh <- fmt.Errorf("Error fetching filename from DB, %v", row)
		}
	}(&filename)
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, "", err
	default:
		return file, filename, nil
	}
}
