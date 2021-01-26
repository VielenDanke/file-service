package model

import "io"

// AWSModel ...
type AWSModel struct {
	File     io.Reader
	FileName string
	IIN      string
	FileLink string
	Metadata map[string]interface{}
}

// GetFile ...
func (am *AWSModel) GetFile() io.Reader {
	return am.File
}

// GetFileName ...
func (am *AWSModel) GetFileName() string {
	return am.FileName
}

// GetIIN ...
func (am *AWSModel) GetIIN() string {
	return am.IIN
}

// GetMetadata ...
func (am *AWSModel) GetMetadata() map[string]interface{} {
	return am.Metadata
}

// GetFileLink ...
func (am *AWSModel) GetFileLink() string {
	return am.FileLink
}
