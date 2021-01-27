package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
	"github.com/vielendanke/file-service/internal/app/fileservice/repository"
)

var awsRepo repository.FileRepository
var mock sqlmock.Sqlmock

func init() {
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()
	if err != nil {
		fmt.Printf("Error during initialization mock: %v", err)
		return
	}
	mockDB := sqlx.NewDb(db, "sqlmock")
	awsRepo = repository.NewAWSFileRepository(mockDB)
}

func TestSaveFile(t *testing.T) {
	testData := "testData"
	fModel := &model.AWSModel{
		FileID:   testData,
		FileLink: testData,
		FileName: testData,
		IIN:      testData,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO FILES").WithArgs(
		testData, testData, testData, testData, testData,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := awsRepo.SaveFile(context.Background(), fModel, testData); err != nil {
		t.Fatalf("Error was not expected while saving file: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestFindFileMetadataByID(t *testing.T) {
	testID := "testID"
	testData := "testData"

	mock.ExpectQuery("SELECT METADATA FROM FILES").WithArgs(testID).WillReturnRows(sqlmock.NewRows([]string{"METADATA"}).AddRow(testData))

	res, err := awsRepo.FindFileMetadataByID(context.Background(), testID)
	if err != nil {
		t.Fatalf("Unexpected error while fetching metadata by ID, %v", err)
	}

	assert.Equal(t, testData, res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}
