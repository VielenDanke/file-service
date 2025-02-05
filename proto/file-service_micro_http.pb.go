// Code generated by protoc-gen-micro
// source: proto/file-service.proto
package pb

import (
	"context"

	micro_client_http "github.com/unistack-org/micro-client-http/v3"
	micro_client "github.com/unistack-org/micro/v3/client"
	micro_server "github.com/unistack-org/micro/v3/server"
)

var (
	_ micro_server.Option
	_ micro_client.Option
)

type fileProcessingService struct {
	c    micro_client.Client
	name string
}

// Micro client stuff

// NewFileProcessingService create new service client
func NewFileProcessingService(name string, c micro_client.Client) FileProcessingService {
	return &fileProcessingService{c: c, name: name}
}

func (c *fileProcessingService) FileProcessing(ctx context.Context, req *FileProcessingRequest, opts ...micro_client.CallOption) (*FileProcessingResponse, error) {
	nopts := append(opts,
		micro_client_http.Method("POST"),
		micro_client_http.Path("/files"),
	)
	rsp := &FileProcessingResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "FileProcessing.FileProcessing", req), rsp, nopts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *fileProcessingService) GetFileMetadata(ctx context.Context, req *GetMetadataRequest, opts ...micro_client.CallOption) (*GetMetadataResponse, error) {
	nopts := append(opts,
		micro_client_http.Method("GET"),
		micro_client_http.Path("/metadata/{metadata_id}"),
	)
	rsp := &GetMetadataResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "FileProcessing.GetFileMetadata", req), rsp, nopts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *fileProcessingService) DownloadFile(ctx context.Context, req *FileDownloadRequest, opts ...micro_client.CallOption) (*FileDownloadResponse, error) {
	nopts := append(opts,
		micro_client_http.Method("GET"),
		micro_client_http.Path("/files/{file_download_id}"),
	)
	rsp := &FileDownloadResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "FileProcessing.DownloadFile", req), rsp, nopts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *fileProcessingService) UpdateFileMetadata(ctx context.Context, req *UpdateMetadataRequest, opts ...micro_client.CallOption) (*UpdateMetadataResponse, error) {
	nopts := append(opts,
		micro_client_http.Method("PUT"),
		micro_client_http.Path("/metadata/{update_metadata_id}"),
	)
	rsp := &UpdateMetadataResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "FileProcessing.UpdateFileMetadata", req), rsp, nopts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// Micro server stuff

type fileProcessingHandler struct {
	FileProcessingHandler
}

func (h *fileProcessingHandler) FileProcessing(ctx context.Context, req *FileProcessingRequest, rsp *FileProcessingResponse) error {
	return h.FileProcessingHandler.FileProcessing(ctx, req, rsp)
}

func (h *fileProcessingHandler) GetFileMetadata(ctx context.Context, req *GetMetadataRequest, rsp *GetMetadataResponse) error {
	return h.FileProcessingHandler.GetFileMetadata(ctx, req, rsp)
}

func (h *fileProcessingHandler) DownloadFile(ctx context.Context, req *FileDownloadRequest, rsp *FileDownloadResponse) error {
	return h.FileProcessingHandler.DownloadFile(ctx, req, rsp)
}

func (h *fileProcessingHandler) UpdateFileMetadata(ctx context.Context, req *UpdateMetadataRequest, rsp *UpdateMetadataResponse) error {
	return h.FileProcessingHandler.UpdateFileMetadata(ctx, req, rsp)
}
