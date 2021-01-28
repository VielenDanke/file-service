package service_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vielendanke/file-service/internal/app/fileservice/mocks"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
	"github.com/vielendanke/file-service/internal/app/fileservice/service"
)

func TestAWSProcessingService_SaveFile(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		FileID:   testData,
		Metadata: metadata,
	}

	mockCodec.On("Marshal", metadata).Return([]byte(testData), nil)
	mockRepo.On("SaveFile", context.Background(), awsModel, testData).Return(nil)

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_SaveFile_ShouldFail_NilFileID(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		Metadata: metadata,
	}

	mockCodec.On("Marshal", metadata).Return([]byte(testData), nil)

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.NotNil(t, err)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_SaveFile_ShouldFail_MarshalReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	metadata := make(map[string]interface{})
	errMsg := "My custom error message"
	awsModel := &model.AWSModel{
		Metadata: metadata,
	}

	mockCodec.On("Marshal", metadata).Return(nil, fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_SaveFile_ShouldFail_RepositoryReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	metadata := make(map[string]interface{})
	errMsg := "my custom error message"
	awsModel := &model.AWSModel{
		FileID:   testData,
		Metadata: metadata,
	}

	mockCodec.On("Marshal", metadata).Return([]byte(testData), nil)
	mockRepo.On("SaveFile", context.Background(), awsModel, testData).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_StoreFile(t *testing.T) {
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		Metadata: metadata,
		File:     bytes.NewBuffer([]byte(testData)),
	}

	mockStore.On(
		"Write",
		context.Background(),
		mock.Anything,
		awsModel.File,
		mock.AnythingOfType("store.WriteOption"),
		mock.AnythingOfType("store.WriteOption"),
	).Return(nil)
	mockStore.On("Exists", context.Background(), mock.Anything, mock.AnythingOfType("store.ExistsOption")).Return(nil)

	awsService := service.NewAWSProcessingService(nil, nil, mockStore)

	err := awsService.StoreFile(context.Background(), awsModel)

	assert.Nil(t, err)
	assert.NotNil(t, awsModel.FileID)
	assert.NotContains(t, awsModel.FileID, "-")

	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_StoreFile_WriterReturnError(t *testing.T) {
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	errMsg := "my custom error message"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		Metadata: metadata,
		File:     bytes.NewBuffer([]byte(testData)),
	}

	mockStore.On(
		"Write",
		context.Background(),
		mock.Anything,
		awsModel.File,
		mock.AnythingOfType("store.WriteOption"),
		mock.AnythingOfType("store.WriteOption"),
	).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, nil, mockStore)

	err := awsService.StoreFile(context.Background(), awsModel)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	assert.NotNil(t, awsModel.FileID)
	assert.NotContains(t, awsModel.FileID, "-")

	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_StoreFile_ExistsReturnError(t *testing.T) {
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	metadata := make(map[string]interface{})
	errMsg := "my custom error message"
	awsModel := &model.AWSModel{
		Metadata: metadata,
		File:     bytes.NewBuffer([]byte(testData)),
	}

	mockStore.On(
		"Write",
		context.Background(),
		mock.Anything,
		awsModel.File,
		mock.AnythingOfType("store.WriteOption"),
		mock.AnythingOfType("store.WriteOption"),
	).Return(nil)
	mockStore.On("Exists", context.Background(), mock.Anything, mock.AnythingOfType("store.ExistsOption")).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, nil, mockStore)

	err := awsService.StoreFile(context.Background(), awsModel)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	assert.NotNil(t, awsModel.FileID)
	assert.NotContains(t, awsModel.FileID, "-")
}

func TestAWSProcessingService_GetFileMetadata(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	testData := "metadata"

	mockRepo.On("FindFileMetadataByID", context.Background(), mock.Anything).Return(testData, nil)

	awsService := service.NewAWSProcessingService(nil, mockRepo, nil)

	res, err := awsService.GetFileMetadata(context.Background(), testData)

	assert.Nil(t, err)
	assert.Equal(t, testData, res)
	mockRepo.AssertExpectations(t)
}

func TestAWSProcessingService_GetFileMetadata_RepositoryReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	testData := "metadata"
	errMsg := "my custom error message"

	mockRepo.On("FindFileMetadataByID", context.Background(), mock.Anything).Return("", fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, mockRepo, nil)

	res, err := awsService.GetFileMetadata(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	assert.Equal(t, "", res)
	mockRepo.AssertExpectations(t)
}

func TestAWSProcessingService_DownloadFile(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	file := []byte{}

	mockStore.On("Read", context.Background(), testData, &file, mock.AnythingOfType("store.ReadOption")).Return(nil)
	mockRepo.On("FindFileNameByID", context.Background(), testData).Return(testData, nil)

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore)

	fileBytes, filename, err := awsService.DownloadFile(context.Background(), testData)

	assert.Nil(t, err)
	assert.NotNil(t, fileBytes)
	assert.Equal(t, testData, filename)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_DownloadFile_StoreReadReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	errMsg := "my custom error message"
	file := []byte{}

	mockStore.On("Read", context.Background(), testData, &file, mock.AnythingOfType("store.ReadOption")).Return(fmt.Errorf(errMsg))
	mockRepo.On("FindFileNameByID", context.Background(), testData).Return(testData, nil)

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore)

	fileBytes, filename, err := awsService.DownloadFile(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Nil(t, fileBytes)
	assert.NotEqual(t, testData, filename)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_DownloadFile_RepositoryFindFileNameByIDReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	errMsg := "my custom error message"
	file := []byte{}

	mockStore.On("Read", context.Background(), testData, &file, mock.AnythingOfType("store.ReadOption")).Return(nil)
	mockRepo.On("FindFileNameByID", context.Background(), testData).Return("", fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore)

	fileBytes, filename, err := awsService.DownloadFile(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Nil(t, fileBytes)
	assert.NotEqual(t, testData, filename)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_DownloadFile_RepositoryFindFileNameByIDReturnError_StoreReadReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	errMsg := "my custom error message"
	file := []byte{}

	mockStore.On("Read", context.Background(), testData, &file, mock.AnythingOfType("store.ReadOption")).Return(fmt.Errorf(errMsg))
	mockRepo.On("FindFileNameByID", context.Background(), testData).Return("", fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore)

	fileBytes, filename, err := awsService.DownloadFile(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Nil(t, fileBytes)
	assert.NotEqual(t, testData, filename)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_UpdateFileMetadata(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "testData"
	testMetadata := make(map[string]interface{})

	mockCodec.On("Marshal", testMetadata).Return([]byte(testData), nil)
	mockRepo.On("UpdateFileMetadataByID", context.Background(), testData, testData).Return(nil)

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.Nil(t, err)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAWSProcessingService_UpdateFileMetadata_MarshalReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "testData"
	errMsg := "my custom error message"
	testMetadata := make(map[string]interface{})

	mockCodec.On("Marshal", testMetadata).Return(nil, fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAWSProcessingService_UpdateFileMetadata_RepositoryUpdateFileMetadataReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "testData"
	errMsg := "my custom error message"
	testMetadata := make(map[string]interface{})

	mockCodec.On("Marshal", testMetadata).Return([]byte(testData), nil)
	mockRepo.On("UpdateFileMetadataByID", context.Background(), testData, testData).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
