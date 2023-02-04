package files

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/config"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/models"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/ports"
)

type UseCase interface {
	UploadFile(ctx context.Context, tenant string, command models.UploadFileCommand) error
	ListBy(ctx context.Context, tenant string, referenceID string) (*[]models.Attachment, error)
	DeleteFile(ctx context.Context, tenant string, fileID string) error
}

type DefaultFilesUseCase struct {
	fileStorage    ports.FileStorage
	fileRepository ports.FileRepository
	idGen          ports.IdGenerator
	expirationTime int
	insecure       bool
}

func NewDefaultFilesUseCase(fileStorage ports.FileStorage, fileRepository ports.FileRepository, idGen ports.IdGenerator, gcloudConfig config.GCloudStorage) (*DefaultFilesUseCase, error) {
	return &DefaultFilesUseCase{fileStorage: fileStorage, expirationTime: gcloudConfig.UrlExpirationTime, insecure: gcloudConfig.Insecure, fileRepository: fileRepository, idGen: idGen}, nil
}

func (useCase *DefaultFilesUseCase) UploadFile(ctx context.Context, tenant string, command models.UploadFileCommand) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(command.File)
	if err != nil {
		return fmt.Errorf("DefaultFilesUseCase - UploadFile: can't read from file; %w", err)
	}
	err = useCase.fileStorage.UploadFile(ctx, tenant, command.FileName, buf.Bytes())
	if err != nil {
		return fmt.Errorf("DefaultFilesUseCase - UploadFile: can't upload file to file storage; %w", err)
	}
	createdFile := models.File{
		ID:          useCase.idGen.MakeId(),
		FileName:    command.FileName,
		ReferenceID: command.ReferenceID,
		CreatedAt:   time.Now().Unix(),
		CreatorId:   command.CreatorId,
	}
	err = useCase.fileRepository.Add(ctx, tenant, &createdFile)
	if err != nil {
		return fmt.Errorf("DefaultFilesUseCase - UploadFile: can't add file to repository; %w", err)
	}
	return nil
}

func (useCase *DefaultFilesUseCase) ListBy(ctx context.Context, tenant string, referenceID string) (*[]models.Attachment, error) {
	uploadedFiles, err := useCase.fileRepository.ListBy(ctx, tenant, referenceID)
	var result []models.Attachment
	if err != nil {
		return nil, fmt.Errorf("DefaultFilesUseCase - ListBy: can't list files by reference id; %w", err)
	}
	for i := range *uploadedFiles {
		url, err := useCase.fileStorage.GetExpiringUrl(tenant, (*uploadedFiles)[i].FileName, useCase.expirationTime, useCase.insecure)
		if err != nil {
			return nil, fmt.Errorf("DefaultFilesUseCase - ListBy: can't get file url; %w", err)
		}
		result = append(result, models.Attachment{
			ID:       (*uploadedFiles)[i].ID,
			FileName: (*uploadedFiles)[i].FileName,
			Url:      url,
		})
	}
	return &result, nil
}

func (useCase *DefaultFilesUseCase) DeleteFile(ctx context.Context, tenant string, fileID string) error {
	file, err := useCase.fileRepository.ReadBy(ctx, tenant, fileID)
	if err != nil {
		return fmt.Errorf("DefaultFilesUseCase - DeleteFile: file doesn't exist; %w", err)
	}
	err = useCase.fileRepository.Delete(ctx, tenant, fileID)
	if err != nil {
		return fmt.Errorf("DefaultFilesUseCase - DeleteFile: can't delete file from repository; %w", err)
	}
	err = useCase.fileStorage.DeleteFile(ctx, tenant, file.FileName)
	if err != nil {
		return fmt.Errorf("DefaultFilesUseCase - DeleteFile: can't delete file from storage; %w", err)
	}
	return nil
}
