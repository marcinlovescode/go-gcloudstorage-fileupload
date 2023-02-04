package ports

import (
	"context"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/models"
)

type FileRepository interface {
	ListBy(ctx context.Context, tenant, referenceID string) (*[]models.File, error)
	ReadBy(ctx context.Context, tenant, fileID string) (*models.File, error)
	Add(ctx context.Context, tenant string, file *models.File) error
	Delete(ctx context.Context, tenant, fileID string) error
}
