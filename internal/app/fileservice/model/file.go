package model

import "io"

// FileModel ...
type FileModel interface {
	GetFile() io.Reader
	GetFileName() string
	GetIIN() string
	GetFileLink() string
	GetMetadata() map[string]interface{}
}
