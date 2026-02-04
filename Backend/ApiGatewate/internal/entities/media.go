package entities

import "io"

type UploadResponse struct {
	Files []FileResp
}

type File struct {
	Filename string
	File     io.ReadSeeker
}

type FileResp struct {
	Filename string
	URL      string
}
