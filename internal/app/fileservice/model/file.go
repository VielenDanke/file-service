package model

import "io"

// FileModel ...
type FileModel interface {
	GetFile() io.Reader
	GetFileID() string
	GetFileName() string
	GetIIN() string
	GetFileLink() string
	GetMetadata() map[string]interface{}
}
