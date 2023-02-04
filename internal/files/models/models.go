package models

import "io"

type UploadFileCommand struct {
	CreatorId   string
	FileName    string
	ReferenceID string
	File        io.Reader
}

type Attachment struct {
	ID       string
	FileName string
	Url      string
}

type File struct {
	ID          string
	FileName    string
	ReferenceID string
	CreatedAt   int64
	CreatorId   string
}
