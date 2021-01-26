package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
)

// AWSFileRepository ...
type AWSFileRepository struct {
	db *sqlx.DB
}

// NewAWSFileRepository ...
func NewAWSFileRepository(db *sqlx.DB) FileRepository {
	return &AWSFileRepository{
		db: db,
	}
}

// SaveFile ...
func (afr *AWSFileRepository) SaveFile(ctx context.Context, f model.FileModel, metadata string) error {
	awsFile := f.(*model.AWSModel)
	tx := afr.db.MustBegin()
	res := tx.MustExecContext(ctx, "INSERT INTO FILES(ID, FILE_NAME, IIN, FILE_LINK, METADATA) VALUES($1, $2, $3, $4, $5)",
		awsFile.GetFileID(), awsFile.GetFileName(), awsFile.GetIIN(), awsFile.GetFileLink(), metadata,
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

// FindFileMetadataByID ...
func (afr *AWSFileRepository) FindFileMetadataByID(ctx context.Context, id string) (string, error) {
	jsonMetadata := ""
	if err := afr.db.QueryRowContext(ctx, "SELECT METADATA FROM FILES WHERE ID=$1", id).Scan(&jsonMetadata); err != nil {
		return "", fmt.Errorf("Error while reading from DB, %v", err)
	}
	return jsonMetadata, nil
}

// FindFileNameByID ...
func (afr *AWSFileRepository) FindFileNameByID(ctx context.Context, id string) (string, error) {
	filename := ""
	if row := afr.db.QueryRowContext(ctx, "SELECT FILE_NAME FROM FILES WHERE ID=$1", id).Scan(&filename); row != nil {
		return "", fmt.Errorf("Error fetching filename from DB, %v", row)
	}
	return filename, nil
}
