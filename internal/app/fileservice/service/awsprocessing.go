package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
)

// AWSProcessingService ...
type AWSProcessingService struct {
	db    *sqlx.DB
	codec codec.Codec
}

// NewAWSProcessingService ...
func NewAWSProcessingService(codec codec.Codec, db *sqlx.DB) FileProcessingService {
	return &AWSProcessingService{
		db:    db,
		codec: codec,
	}
}

// StoreFile ...
func (aps *AWSProcessingService) StoreFile(f model.FileModel) error {
	return nil
}

// SaveFileData ...
func (aps *AWSProcessingService) SaveFileData(f model.FileModel) (string, error) {
	awsFile := f.(*model.AWSModel)
	jsonMetadata, err := aps.codec.Marshal(awsFile.GetMetadata())
	if err != nil {
		return "", fmt.Errorf("Error while marshalling metadata, %v", err)
	}
	fileID := uuid.New().String()
	tx := aps.db.MustBegin()
	tx.MustExec("INSERT INTO FILES(ID, FILE_NAME, IIN, FILE_LINK, METADATA) VALUES($1, $2, $3, $4, $5)",
		fileID, awsFile.GetFileName(), awsFile.GetIIN(), awsFile.GetFileLink(), string(jsonMetadata),
	)
	tx.Commit()
	return fileID, nil
}

// GetFileMetadata ...
func (aps *AWSProcessingService) GetFileMetadata(id string) (string, error) {
	jsonMetadata := ""
	if err := aps.db.QueryRow("SELECT METADATA FROM FILES WHERE ID=$1", id).Scan(&jsonMetadata); err != nil {
		return "", fmt.Errorf("Error while reading from DB, %v", err)
	}
	return jsonMetadata, nil
}
