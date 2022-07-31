package files



import (
	"context"
	"api/v1/video/models"
)

type Repository interface {
	CreateFile(ctx context.Context, file *models.File) error
	UpdateFileProgressById(ctx context.Context, id string, progress int) error
	UpdateFileProcessingStatusById(ctx context.Context, id string, processingStatus string) error
	UpdateFileIsProcessingById(ctx context.Context, id string, isProcessing bool) error
	GetFiles(ctx context.Context) ([]*models.File,error)
	GetFileById(ctx context.Context,id string) (*models.File,error)
	DeleteFile(ctx context.Context, id string) error
}