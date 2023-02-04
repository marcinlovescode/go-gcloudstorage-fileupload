package adapters

import (
	"context"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/models"
)

type InMemoryFilesRepository struct {
	files map[string][]models.File
}

func NewInMemoryFilesRepository(_ context.Context) (*InMemoryFilesRepository, error) {
	return &InMemoryFilesRepository{
		files: make(map[string][]models.File),
	}, nil
}

func (repository *InMemoryFilesRepository) ListBy(_ context.Context, tenant, referenceID string) (*[]models.File, error) {
	data, ok := repository.files[tenant]
	if !ok {
		return &[]models.File{}, nil
	}
	var result []models.File

	for i := range data {
		if data[i].ReferenceID == referenceID {
			result = append(result, data[i])
		}
	}
	return &result, nil
}

func (repository *InMemoryFilesRepository) ReadBy(_ context.Context, tenant, fileID string) (*models.File, error) {
	data, ok := repository.files[tenant]
	if !ok {
		return &models.File{}, nil
	}

	for i := range data {
		if data[i].ID == fileID {
			return &data[i], nil
		}
	}
	return nil, nil
}

func (repository *InMemoryFilesRepository) Add(_ context.Context, tenant string, file *models.File) error {
	data, ok := repository.files[tenant]
	if !ok {
		repository.files[tenant] = []models.File{}
		data = repository.files[tenant]
	}
	repository.files[tenant] = append(data, *file)
	return nil
}

func (repository *InMemoryFilesRepository) Delete(_ context.Context, tenant, fileID string) error {
	data, ok := repository.files[tenant]
	if !ok {
		return nil
	}
	for i := range data {
		if data[i].ID == fileID {
			repository.files[tenant] = append(data[:i], data[i+1:]...)
			break
		}
	}
	return nil
}
