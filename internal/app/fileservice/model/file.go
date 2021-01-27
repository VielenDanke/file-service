package model

import "io"

// FileModel ...
type FileModel interface {
	GetFile() io.Reader
	GetFileID() string
	GetFileName() string
	GetDocClass() string
	GetDocType() string
	GetDocNum() string
	GetMetadata() map[string]interface{}
}
