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

func TestAWSProcessingService_SaveFileMetadata(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		FileID:   testData,
		Metadata: metadata,
	}

	mockRepo.On("CheckIfExists", mock.Anything, awsModel).Return(nil)
	mockCodec.On("Marshal", metadata).Return([]byte(testData), nil)
	mockRepo.On("SaveFileMetadata", context.Background(), awsModel, testData).Return(nil)

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.Nil(t, err)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_SaveFileMetadata_FileMetadataAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	errMsg := "Error"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		FileID:   testData,
		Metadata: metadata,
	}

	mockRepo.On("CheckIfExists", mock.Anything, awsModel).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.NotNil(t, err)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_SaveFileMetadata_CodecMarshalReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	metadata := make(map[string]interface{})
	errMsg := "My custom error message"
	awsModel := &model.AWSModel{
		Metadata: metadata,
	}

	mockRepo.On("CheckIfExists", mock.Anything, awsModel).Return(nil)
	mockCodec.On("Marshal", metadata).Return(nil, fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil, nil)

	err := awsService.SaveFileData(context.Background(), awsModel)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)

	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
}

func TestAWSProcessingService_SaveFileMetadata_RepositorySaveFileReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	metadata := make(map[string]interface{})
	errMsg := "my custom error message"
	awsModel := &model.AWSModel{
		FileID:   testData,
		Metadata: metadata,
	}

	mockRepo.On("CheckIfExists", mock.Anything, awsModel).Return(nil)
	mockCodec.On("Marshal", metadata).Return([]byte(testData), nil)
	mockRepo.On("SaveFileMetadata", context.Background(), awsModel, testData).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, nil, nil)

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
		FileID:   testData,
	}

	mockStore.On(
		"Write",
		context.Background(),
		mock.Anything,
		awsModel.File,
		mock.AnythingOfType("store.WriteOption"),
		mock.AnythingOfType("store.WriteOption"),
	).Return(nil)

	awsService := service.NewAWSProcessingService(nil, nil, mockStore, mockStore)

	err := awsService.StoreFile(context.Background(), awsModel)

	assert.Nil(t, err)

	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_StoreFile_EmptyFileID(t *testing.T) {
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		Metadata: metadata,
		File:     bytes.NewBuffer([]byte(testData)),
	}

	awsService := service.NewAWSProcessingService(nil, nil, mockStore, mockStore)

	err := awsService.StoreFile(context.Background(), awsModel)

	assert.NotNil(t, err)

	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_StoreFile_StoreWriterReturnError(t *testing.T) {
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	errMsg := "my custom error message"
	metadata := make(map[string]interface{})
	awsModel := &model.AWSModel{
		Metadata: metadata,
		File:     bytes.NewBuffer([]byte(testData)),
		FileID:   testData,
	}

	mockStore.On(
		"Write",
		context.Background(),
		mock.Anything,
		awsModel.File,
		mock.AnythingOfType("store.WriteOption"),
		mock.AnythingOfType("store.WriteOption"),
	).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, nil, mockStore, mockStore)

	err := awsService.StoreFile(context.Background(), awsModel)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)

	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_GetFileMetadata(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	testMetadata := make(map[string]string)
	testMetadata["metadata"] = testData

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(nil)
	mockRepo.On("FindFileMetadataByID", context.Background(), mock.Anything).Return(testMetadata, nil)
	mockCodec.On("Unmarshal", []byte(testData), mock.Anything).Return(nil)

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	res, err := awsService.GetFileMetadata(context.Background(), testData)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_GetFileMetadata_StoreExistsReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	errMsg := "error"
	testMetadata := make(map[string]string)
	testMetadata["metadata"] = testData

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	res, err := awsService.GetFileMetadata(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_GetFileMetadata_RepositoryReturnFindFileMetadataByIDError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	testData := "metadata"
	errMsg := "my custom error message"

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(nil)
	mockRepo.On("FindFileMetadataByID", context.Background(), mock.Anything).Return(nil, fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore, nil)

	res, err := awsService.GetFileMetadata(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_GetFileMetadata_CodecUnmarshalReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	mockCodec := new(mocks.MockCodec)
	testData := "metadata"
	errMsg := "my custom error message"

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(nil)
	mockRepo.On("FindFileMetadataByID", context.Background(), mock.Anything).Return(make(map[string]string), nil)
	mockCodec.On("Unmarshal", mock.Anything, mock.Anything).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	res, err := awsService.GetFileMetadata(context.Background(), testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_DownloadFile(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	file := []byte{}

	mockStore.On("Read", context.Background(), testData, &file, mock.AnythingOfType("store.ReadOption")).Return(nil)
	mockRepo.On("FindFileNameByID", context.Background(), testData).Return(testData, nil)

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore, mockStore)

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

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore, mockStore)

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

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore, mockStore)

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

	awsService := service.NewAWSProcessingService(nil, mockRepo, mockStore, mockStore)

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
	mockStore := new(mocks.MockStore)
	testData := "testData"
	testMetadata := make(map[string]interface{})

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(nil)
	mockCodec.On("Marshal", testMetadata).Return([]byte(testData), nil)
	mockRepo.On("UpdateFileMetadataByID", context.Background(), testData, testData).Return(nil)

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.Nil(t, err)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_UpdateFileMetadata_StoreExistsReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	errMsg := "error"
	testMetadata := make(map[string]interface{})

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.NotNil(t, err)
	mockStore.AssertExpectations(t)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAWSProcessingService_UpdateFileMetadata_MarshalReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	errMsg := "my custom error message"
	testMetadata := make(map[string]interface{})

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(nil)
	mockCodec.On("Marshal", testMetadata).Return(nil, fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_UpdateFileMetadata_RepositoryUpdateFileMetadataReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	mockCodec := new(mocks.MockCodec)
	mockStore := new(mocks.MockStore)
	testData := "testData"
	errMsg := "my custom error message"
	testMetadata := make(map[string]interface{})

	mockStore.On("Exists", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("store.ExistsOption")).Return(nil)
	mockCodec.On("Marshal", testMetadata).Return([]byte(testData), nil)
	mockRepo.On("UpdateFileMetadataByID", context.Background(), testData, testData).Return(fmt.Errorf(errMsg))

	awsService := service.NewAWSProcessingService(mockCodec, mockRepo, mockStore, nil)

	err := awsService.UpdateFileMetadata(context.Background(), testMetadata, testData)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMsg)
	mockCodec.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestAWSProcessingService_DeleteMetadataByID(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	testData := "testData"

	mockRepo.On("DeleteMetadataByID", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	awsService := service.NewAWSProcessingService(nil, mockRepo, nil, nil)

	err := awsService.DeleteMetadataByID(context.Background(), testData)

	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAWSProcessingService_DeleteMetadataByID_RepositoryDeleteByIDReturnError(t *testing.T) {
	mockRepo := new(mocks.FileRepository)
	testData := "testData"

	mockRepo.On("DeleteMetadataByID", mock.Anything, mock.AnythingOfType("string")).Return(fmt.Errorf(testData))

	awsService := service.NewAWSProcessingService(nil, mockRepo, nil, nil)

	err := awsService.DeleteMetadataByID(context.Background(), testData)

	assert.NotNil(t, err)
	mockRepo.AssertExpectations(t)
}
