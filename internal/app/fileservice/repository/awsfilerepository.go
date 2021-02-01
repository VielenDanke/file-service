package repository

import (
	"context"
	"database/sql"
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

// CheckIfExists ...
func (afr *AWSFileRepository) CheckIfExists(ctx context.Context, fields ...string) error {
	rows := afr.db.QueryRowContext(ctx, "SELECT ID FROM FILES WHERE DOC_CLASS=$1 AND DOC_TYPE=$2 AND DOC_NUM=$3", fields[0], fields[1], fields[2])
	if rows.Scan() == sql.ErrNoRows {
		return nil
	}
	return fmt.Errorf("File alredy exists in the system, to update - use update enpdoint")
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
func (afr *AWSFileRepository) FindFileMetadataByID(ctx context.Context, id string) (map[string]string, error) {
	jsonMetadata := ""
	docType := ""
	docClass := ""
	docNum := ""
	if err := afr.db.QueryRowContext(
		ctx,
		"SELECT DOC_CLASS, DOC_TYPE, DOC_NUM, METADATA FROM FILES WHERE ID=$1",
		id,
	).Scan(&docClass, &docType, &docNum, &jsonMetadata); err != nil {
		return nil, fmt.Errorf("Error while reading from DB, %v", err)
	}
	properties := make(map[string]string)
	properties["type"] = docType
	properties["class"] = docClass
	properties["number"] = docNum
	properties["metadata"] = jsonMetadata
	return properties, nil
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
