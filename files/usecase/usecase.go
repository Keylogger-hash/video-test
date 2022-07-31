package usecase


import (
	"context"
	"api/v1/video/files"
	"api/v1/video/models"
	"github.com/google/uuid"

)


type FileUseCase struct {
	fileRepo files.Repository
}


func NewFileUseCase(fileRepo files.Repository) *FileUseCase {
	return &FileUseCase{
		fileRepo: fileRepo,
	}
}

func (f FileUseCase) CreateFile(ctx context.Context,path string, filename string) error {
	uuid := uuid.New().String()
	fm := &models.File{
		Id: uuid,
		Path: path,
		Filename: filename,
		Processing: false,
		Progress: 0,
		ProcessingStatus: "UPLOAD",
	}
	return f.fileRepo.CreateFile(ctx,fm)
}

func (f FileUseCase)  UpdateFileProgressById(ctx context.Context,id string, progress int) error {
	return f.fileRepo.UpdateFileProgressById(ctx,id,progress)
}

func (f FileUseCase)  UpdateFileProcessingStatusById(ctx context.Context,id string, processingStatus string) error {
	return f.fileRepo.UpdateFileProcessingStatusById(ctx,id,processingStatus)
}

func (f FileUseCase)  UpdateFileIsProcessingById(ctx context.Context,id string, isProcessing bool) error {
	return f.fileRepo.UpdateFileIsProcessingById(ctx,id,isProcessing)
}