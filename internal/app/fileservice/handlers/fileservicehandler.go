package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/vielendanke/file-service/internal/app/fileservice/model"
	"github.com/vielendanke/file-service/internal/app/fileservice/service"
	"github.com/vielendanke/file-service/internal/app/fileservice/validations"
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
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, fmt.Sprintf("Failed to parse request file, %v", err))
		return
	}
	jsonBody := r.FormValue("body")
	if jsonBody == "" {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, fmt.Sprintf("Bad request, body is empty"))
		return
	}
	properties := make(map[string]interface{})
	if err := fh.codec.Unmarshal([]byte(jsonBody), &properties); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, fmt.Sprintf("Error unmarshalling request, %v", err))
		return
	}
	if err := validations.ValidateJSONDocumentRequest(properties); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, err.Error())
		return
	}
	awsFile := &model.AWSModel{
		File:     file,
		FileName: header.Filename,
		DocClass: properties["class"].(string),
		DocType:  properties["type"].(string),
		DocNum:   properties["number"].(string),
		Metadata: properties,
	}
	if saveErr := fh.service.SaveFileData(r.Context(), awsFile); saveErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, fmt.Sprintf("Error saving file metadata to DB, %v", saveErr))
		return
	}
	if err := fh.service.StoreFile(r.Context(), awsFile); err != nil {
		if deleteErr := fh.service.DeleteMetadataByID(r.Context(), awsFile.GetFileID()); deleteErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fh.codec.Write(w, nil, fmt.Sprintf("Error storing file %v and delete metadata %v", err, deleteErr))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, fmt.Sprintf("Error storing file to s3, %v", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	fh.codec.Write(w, nil, awsFile.GetFileID())
}

// GetFileMetadata ...
func (fh *FileServiceHandler) GetFileMetadata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["metadata_id"]
	metadata, err := fh.service.GetFileMetadata(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fh.codec.Write(w, nil, metadata)
}

// DownloadFile ...
func (fh *FileServiceHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["file_download_id"]
	file, filename, err := fh.service.DownloadFile(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(file)
}

// UpdateFileMetadata ...
func (fh *FileServiceHandler) UpdateFileMetadata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["update_metadata_id"]
	properties := make(map[string]interface{})
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, "Error, body is nil")
		return
	}
	err := fh.codec.ReadBody(r.Body, &properties)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fh.codec.Write(w, nil, fmt.Sprintf("Error reading body, %v", err))
		return
	}
	defer r.Body.Close()
	if uErr := fh.service.UpdateFileMetadata(r.Context(), properties, id); uErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fh.codec.Write(w, nil, fmt.Sprintf("Error updating metadata, %v", uErr))
		return
	}
	w.WriteHeader(http.StatusOK)
	fh.codec.Write(w, nil, properties)
}
