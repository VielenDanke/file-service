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

func setupDB() {
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

func TestDeleteMetadataByID(t *testing.T) {
	setupDB()
	testData := "testData"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testData).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := awsRepo.DeleteMetadataByID(context.Background(), testData)

	assert.Nil(t, err)
	mock.ExpectationsWereMet()
}

func TestDeleteMetadataByID_NoChangesDetected(t *testing.T) {
	setupDB()
	testData := "testData"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testData).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := awsRepo.DeleteMetadataByID(context.Background(), testData)

	assert.NotNil(t, err)
	mock.ExpectationsWereMet()
}

func TestDeleteMetadataByID_QueryReturnError(t *testing.T) {
	setupDB()
	testData := "testData"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(testData).WillReturnError(fmt.Errorf(testData))
	mock.ExpectRollback()

	err := awsRepo.DeleteMetadataByID(context.Background(), testData)

	assert.NotNil(t, err)
	mock.ExpectationsWereMet()
}

func TestCheckIfExists(t *testing.T) {
	setupDB()
	testData := "testData"
	awsModel := &model.AWSModel{
		DocClass: testData,
		DocType:  testData,
		DocNum:   testData,
	}

	mock.ExpectQuery("SELECT").WithArgs(testData, testData, testData).WillReturnRows(sqlmock.NewRows([]string{}))

	err := awsRepo.CheckIfExists(context.Background(), awsModel)

	assert.Nil(t, err)
	mock.ExpectationsWereMet()
}

func TestCheckIfExists_NoUniqueResult(t *testing.T) {
	setupDB()
	testData := "testData"
	awsModel := &model.AWSModel{
		DocClass: testData,
		DocType:  testData,
		DocNum:   testData,
	}

	mock.ExpectQuery("SELECT").WithArgs(testData, testData, testData).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(testData))

	err := awsRepo.CheckIfExists(context.Background(), awsModel)

	assert.NotNil(t, err)
	mock.ExpectationsWereMet()
}

func TestSaveFileMetadata(t *testing.T) {
	setupDB()
	testData := "testData"
	fModel := &model.AWSModel{
		FileID:   testData,
		FileName: testData,
		DocClass: testData,
		DocType:  testData,
		DocNum:   testData,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO FILES").WithArgs(
		testData, testData, testData, testData, testData, testData,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := awsRepo.SaveFileMetadata(context.Background(), fModel, testData); err != nil {
		t.Fatalf("Error was not expected while saving file: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestSaveFileMetadata_InsertNotHappened(t *testing.T) {
	setupDB()
	testData := "testData"
	fModel := &model.AWSModel{
		FileID:   testData,
		FileName: testData,
		DocClass: testData,
		DocType:  testData,
		DocNum:   testData,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO FILES").WithArgs(
		testData, testData, testData, testData, testData, testData,
	).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := awsRepo.SaveFileMetadata(context.Background(), fModel, testData)

	assert.NotNil(t, err)

	if expectedErr := mock.ExpectationsWereMet(); expectedErr != nil {
		t.Fatalf("Results are not expected: %v", expectedErr)
	}
}

func TestSaveFileMetadata_ReturnError(t *testing.T) {
	setupDB()
	testData := "testData"
	errMessage := "My custom error message"
	fModel := &model.AWSModel{
		FileID:   testData,
		FileName: testData,
		DocClass: testData,
		DocType:  testData,
		DocNum:   testData,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO FILES").WithArgs(
		testData, testData, testData, testData, testData, testData,
	).WillReturnError(fmt.Errorf(errMessage))
	mock.ExpectRollback()

	err := awsRepo.SaveFileMetadata(context.Background(), fModel, testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMessage)

	if expectedErr := mock.ExpectationsWereMet(); expectedErr != nil {
		t.Fatalf("Results are not expected: %v", expectedErr)
	}
}

func TestFindFileMetadataByID(t *testing.T) {
	setupDB()
	testID := "testID"
	testData := "testData"

	mock.ExpectQuery("SELECT").WithArgs(testID).WillReturnRows(
		sqlmock.NewRows(
			[]string{"DOC_CLASS", "DOC_TYPE", "DOC_NUM", "METADATA"},
		).AddRow(testData, testData, testData, testData))

	res, err := awsRepo.FindFileMetadataByID(context.Background(), testID)
	if err != nil {
		t.Fatalf("Unexpected error while fetching metadata by ID, %v", err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestFindFileMetadataByID_NoRowsFound(t *testing.T) {
	setupDB()
	testID := "testID"

	mock.ExpectQuery("SELECT").WithArgs(testID).WillReturnError(sql.ErrNoRows)

	_, err := awsRepo.FindFileMetadataByID(context.Background(), testID)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestFindFileNameByID(t *testing.T) {
	setupDB()
	testID := "testID"
	testData := "testData"

	mock.ExpectQuery("SELECT FILE_NAME FROM FILES").WithArgs(testID).WillReturnRows(sqlmock.NewRows([]string{"FILE_NAME"}).AddRow(testData))

	res, err := awsRepo.FindFileNameByID(context.Background(), testID)
	if err != nil {
		t.Fatalf("Unexpected error while fetching metadata by ID, %v", err)
	}

	assert.Equal(t, testData, res)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestFindFileNameByID_NoRowsFound(t *testing.T) {
	setupDB()
	testID := "testID"

	mock.ExpectQuery("SELECT FILE_NAME FROM FILES").WithArgs(testID).WillReturnError(sql.ErrNoRows)

	_, err := awsRepo.FindFileNameByID(context.Background(), testID)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestUpdateFileMetadataByID(t *testing.T) {
	setupDB()
	testID := "testID"
	testData := "testData"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE FILES").WithArgs(testData, testID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := awsRepo.UpdateFileMetadataByID(context.Background(), testData, testID); err != nil {
		t.Fatalf("Unexpected error while updating metadata by ID, %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Results are not expected: %v", err)
	}
}

func TestUpdateFileMetadataByID_NoRowsAffected(t *testing.T) {
	setupDB()
	testID := "testID"
	testData := "testData"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE FILES").WithArgs(testData, testID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := awsRepo.UpdateFileMetadataByID(context.Background(), testData, testID)

	assert.NotNil(t, err)

	if expectedErr := mock.ExpectationsWereMet(); expectedErr != nil {
		t.Fatalf("Results are not expected: %v", expectedErr)
	}
}

func TestUpdateFileMetadataByID_ReturnError(t *testing.T) {
	setupDB()
	testID := "testID"
	testData := "testData"
	errMessage := "My custom error message"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE FILES").WithArgs(testData, testID).WillReturnError(fmt.Errorf(errMessage))
	mock.ExpectRollback()

	err := awsRepo.UpdateFileMetadataByID(context.Background(), testData, testID)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMessage)

	if expectedErr := mock.ExpectationsWereMet(); expectedErr != nil {
		t.Fatalf("Results are not expected: %v", expectedErr)
	}
}
