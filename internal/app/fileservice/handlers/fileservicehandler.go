package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
	"github.com/vielendanke/file-service/internal/app/fileservice/service"
)

// FileServiceHandler ...
type FileServiceHandler struct {
	codec   codec.Codec
	service service.FileProcessingService
}

// NewFileServiceHandler ...
func NewFileServiceHandler(srv service.FileProcessingService, codec codec.Codec) *FileServiceHandler {
	return &FileServiceHandler{
		service: srv,
		codec:   codec,
	}
}

// FileProcessing ...
func (fh *FileServiceHandler) FileProcessing(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, fmt.Errorf("Failed to parse request file, %v", err))
		return
	}
	queryValues := r.URL.Query()
	if len(queryValues["iin"]) == 0 || len(queryValues["filename"]) == 0 ||
		len(queryValues["iin"][0]) < 12 || len(queryValues["filename"][0]) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, fmt.Errorf("IIN or Filename cannot be empty"))
		return
	}
	metadata := make(map[string]interface{})
	for k, v := range queryValues {
		if k == "iin" || k == "filename" {
			continue
		}
		if len(v) == 1 {
			metadata[k] = v[0]
			continue
		}
		metadata[k] = v
	}
	awsFile := &model.AWSModel{
		IIN: queryValues["iin"][0], FileName: queryValues["filename"][0], Metadata: metadata, File: file,
	}
	if err := fh.service.StoreFile(r.Context(), awsFile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, fmt.Errorf("Error storing file to s3, %v", err))
		return
	}
	if saveErr := fh.service.SaveFileData(r.Context(), awsFile); saveErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, fmt.Errorf("Error saving file metadata to DB, %v", saveErr))
		return
	}
	w.WriteHeader(http.StatusOK)
	fh.codec.Write(w, nil, awsFile.GetFileID())
}

// GetFileMetadata ...
func (fh *FileServiceHandler) GetFileMetadata(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["file_metadata_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metadata, err := fh.service.GetFileMetadata(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metadata))
}

// DownloadFile ...
func (fh *FileServiceHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["file_download_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, filename, err := fh.service.DownloadFile(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(file)
}
