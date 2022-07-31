package models


type File struct {
	Id string
	Path string
	Filename string
	Processing bool
	Progress int
	ProcessingStatus string // UPLOAD, START, FAILURE, WAITING
}
