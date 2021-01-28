package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/vielendanke/file-service/configs"
	"github.com/vielendanke/file-service/internal/app/fileservice/handlers"
	"github.com/vielendanke/file-service/internal/app/fileservice/middlewares"
	"github.com/vielendanke/file-service/internal/app/fileservice/mocks"
	"github.com/vielendanke/file-service/internal/app/fileservice/service"
	pb "github.com/vielendanke/file-service/proto"
)

func prepareMultipartRequest(jsonBody map[string]interface{}) (*multipart.Writer, *bytes.Buffer, error) {
	testData := "testData"
	bodyBytes := []byte{}
	var marshalErr error

	if jsonBody != nil {
		bodyBytes, marshalErr = json.Marshal(jsonBody)
		if marshalErr != nil {
			return nil, nil, fmt.Errorf("Unable to marshal body, %v", marshalErr)
		}
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("body", string(bodyBytes))
	part, formErr := writer.CreateFormFile("file", "file.txt")
	if formErr != nil {
		return nil, nil, fmt.Errorf("Unable to create form file, %v", formErr)
	}
	part.Write([]byte(testData))
	writer.Close()

	return writer, body, nil
}

func preparteRouter(service service.FileProcessingService, codec codec.Codec) (*mux.Router, error) {
	handler := handlers.NewFileServiceHandler(service, codec)
	router := mux.NewRouter()
	router.Use(middlewares.NewContentTypeMiddleware("application/json").ContentTypeMiddleware)
	endpoints := pb.NewFileProcessingServiceEndpoints()

	if endpointsErr := configs.ConfigureHandlerToEndpoints(router, handler, endpoints); endpointsErr != nil {
		return nil, fmt.Errorf("Unable to configure endpoints, %v", endpointsErr)
	}
	return router, nil
}

func TestFileServiceHandler_FileProcessing(t *testing.T) {
	mockService := new(mocks.FileProcessingService)
	jsonBody := make(map[string]interface{})
	jsonBody["class"] = "class"
	jsonBody["type"] = "type"
	jsonBody["number"] = "number"
	jsonBody["iin"] = "iin"

	writer, body, err := prepareMultipartRequest(jsonBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	router, routerErr := preparteRouter(mockService, jsoncodec.NewCodec())
	if routerErr != nil {
		t.Fatal(err.Error())
	}

	rec := httptest.NewRecorder()

	req, reqErr := http.NewRequest(http.MethodPost, "/files", body)

	if reqErr != nil {
		t.Fatalf("Error creating http request, %v", reqErr)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mockService.On("StoreFile", mock.Anything, mock.Anything).Return(nil)
	mockService.On("SaveFileData", mock.Anything, mock.Anything).Return(nil)

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header()["Content-Type"][0])
	mockService.AssertExpectations(t)
}

func TestFileServiceHandler_FileProcessing_FileServiceStoreFileReturnError(t *testing.T) {
	mockService := new(mocks.FileProcessingService)
	errMsg := "error message"
	jsonBody := make(map[string]interface{})
	jsonBody["class"] = "class"
	jsonBody["type"] = "type"
	jsonBody["number"] = "number"
	jsonBody["iin"] = "iin"

	writer, body, err := prepareMultipartRequest(jsonBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	router, routerErr := preparteRouter(mockService, jsoncodec.NewCodec())
	if routerErr != nil {
		t.Fatal(routerErr.Error())
	}
	rec := httptest.NewRecorder()

	req, reqErr := http.NewRequest(http.MethodPost, "/files", body)
	if reqErr != nil {
		t.Fatalf("Error creating http request, %v", reqErr)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mockService.On("StoreFile", mock.Anything, mock.Anything).Return(fmt.Errorf(errMsg))

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header()["Content-Type"][0])
	mockService.AssertExpectations(t)
}

func TestFileServiceHandler_FileProcessing_FileServiceSaveFileDataReturnError(t *testing.T) {
	mockService := new(mocks.FileProcessingService)
	errMsg := "error message"
	jsonBody := make(map[string]interface{})
	jsonBody["class"] = "class"
	jsonBody["type"] = "type"
	jsonBody["number"] = "number"
	jsonBody["iin"] = "iin"

	writer, body, err := prepareMultipartRequest(jsonBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	router, routerErr := preparteRouter(mockService, jsoncodec.NewCodec())
	if routerErr != nil {
		t.Fatal(routerErr.Error())
	}
	rec := httptest.NewRecorder()

	req, reqErr := http.NewRequest(http.MethodPost, "/files", body)
	if reqErr != nil {
		t.Fatalf("Error creating http request, %v", reqErr)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mockService.On("StoreFile", mock.Anything, mock.Anything).Return(nil)
	mockService.On("SaveFileData", mock.Anything, mock.Anything).Return(fmt.Errorf(errMsg))

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header()["Content-Type"][0])
	mockService.AssertExpectations(t)
}

func TestFileServiceHandler_FileProcessing_NoFile(t *testing.T) {
	mockService := new(mocks.FileProcessingService)

	router, routerErr := preparteRouter(mockService, jsoncodec.NewCodec())
	if routerErr != nil {
		t.Fatal(routerErr.Error())
	}
	rec := httptest.NewRecorder()

	req, reqErr := http.NewRequest(http.MethodPost, "/files", nil)
	if reqErr != nil {
		t.Fatalf("Error creating http request, %v", reqErr)
	}
	req.Header.Set("Content-Type", "multipart/form-data")

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header()["Content-Type"][0])
	mockService.AssertExpectations(t)
}

func TestFileServiceHandler_FileProcessing_NoBody(t *testing.T) {
	mockService := new(mocks.FileProcessingService)

	writer, body, err := prepareMultipartRequest(nil)
	if err != nil {
		t.Fatal(err)
	}

	router, routerErr := preparteRouter(mockService, jsoncodec.NewCodec())
	if routerErr != nil {
		t.Fatal(routerErr.Error())
	}
	rec := httptest.NewRecorder()

	req, reqErr := http.NewRequest(http.MethodPost, "/files", body)
	if reqErr != nil {
		t.Fatalf("Error creating http request, %v", reqErr)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header()["Content-Type"][0])
	mockService.AssertExpectations(t)
}

func TestFileServiceHandler_FileProcessing_NotValidBody(t *testing.T) {
	mockService := new(mocks.FileProcessingService)
	jsonBody := make(map[string]interface{})
	jsonBody["class"] = "class"
	jsonBody["type"] = "type"

	writer, body, err := prepareMultipartRequest(jsonBody)
	if err != nil {
		t.Fatal(err)
	}

	router, routerErr := preparteRouter(mockService, jsoncodec.NewCodec())
	if routerErr != nil {
		t.Fatal(routerErr.Error())
	}
	rec := httptest.NewRecorder()

	req, reqErr := http.NewRequest(http.MethodPost, "/files", body)
	if reqErr != nil {
		t.Fatalf("Error creating http request, %v", reqErr)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	assert.Equal(t, "application/json", rec.Header()["Content-Type"][0])
	mockService.AssertExpectations(t)
}
