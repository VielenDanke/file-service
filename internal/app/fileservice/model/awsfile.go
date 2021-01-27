package model

import "io"

// AWSModel ...
type AWSModel struct {
	File     io.Reader
	FileID   string
	FileName string
	DocType  string
	DocNum   string
	DocClass string
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

// GetMetadata ...
func (am *AWSModel) GetMetadata() map[string]interface{} {
	return am.Metadata
}

// GetFileID ...
func (am *AWSModel) GetFileID() string {
	return am.FileID
}

// GetDocType ...
func (am *AWSModel) GetDocType() string {
	return am.DocType
}

// GetDocNum ...
func (am *AWSModel) GetDocNum() string {
	return am.DocNum
}

// GetDocClass ...
func (am *AWSModel) GetDocClass() string {
	return am.DocClass
}
