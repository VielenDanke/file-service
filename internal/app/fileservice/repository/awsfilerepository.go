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
	res, err := tx.ExecContext(ctx, "INSERT INTO FILES(ID, FILE_NAME, DOC_CLASS, DOC_TYPE, DOC_NUM, METADATA) VALUES($1, $2, $3, $4, $5, $6)",
		awsFile.GetFileID(), awsFile.GetFileName(), awsFile.GetDocClass(), awsFile.GetDocType(), awsFile.GetDocNum(), metadata,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error inserting new file to DB, %v", err)
	}
	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		tx.Rollback()
		return fmt.Errorf("Error during rows inserter, %v", rowsErr)
	}
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("No insertions found, %d", rows)
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

// UpdateFileMetadataByID ...
func (afr *AWSFileRepository) UpdateFileMetadataByID(ctx context.Context, metadata, id string) error {
	tx := afr.db.MustBegin()
	res, err := tx.ExecContext(ctx, "UPDATE FILES SET METADATA=$1 WHERE ID=$2", metadata, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error updating metadata, %v", err)
	}
	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		tx.Rollback()
		return fmt.Errorf("Error inserting data in DB, %v", rowsErr)
	}
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("No insertions found")
	}
	tx.Commit()
	return nil
}