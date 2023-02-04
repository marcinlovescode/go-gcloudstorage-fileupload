package gcloudstorage

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCloudAuthCredentials struct {
	AccessID   string
	PrivateKey []byte
}

type gCloudStorageService struct {
	client      *storage.Client
	projectName string
	bucketName  string
	credentials *GCloudAuthCredentials
}

func NewGCloudStorageService(client *storage.Client, projectName, bucketName string, credentials *GCloudAuthCredentials) *gCloudStorageService {
	return &gCloudStorageService{
		client:      client,
		projectName: projectName,
		bucketName:  bucketName,
		credentials: credentials,
	}
}

func (storageService *gCloudStorageService) UploadFile(context context.Context, tenant, fileName string, file []byte) error {
	storageObject := storageService.client.Bucket(storageService.bucketName).Object(fmt.Sprintf("%s/%s", tenant, fileName))
	storageObject = storageObject.If(storage.Conditions{DoesNotExist: true})
	writer := storageObject.NewWriter(context)
	fileExtension := filepath.Ext(fileName)
	writer.ContentType = mime.TypeByExtension(fileExtension)
	writer.ContentDisposition = fmt.Sprintf("attachment; filename=\"%s\"", fileName)
	if _, err := writer.Write(file); err != nil {
		return fmt.Errorf("gcloudstorage - UploadFile: can't upload file; %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("gcloudstorage - UploadFile: can't close writer; %w", err)
	}
	return nil
}

func (storageService *gCloudStorageService) DeleteFile(context context.Context, tenant, fileName string) error {
	storageObject := storageService.client.Bucket(storageService.bucketName).Object(fmt.Sprintf("%s/%s", tenant, fileName))
	err := storageObject.Delete(context)
	if err != nil {
		return fmt.Errorf("gcloudstorage - DeleteFile: can't remove file; %w", err)
	}
	return nil
}

func (storageService *gCloudStorageService) FileExists(context context.Context, tenant, fileName string) (bool, error) {
	query := &storage.Query{
		Prefix: fmt.Sprintf("%s/%s", tenant, fileName),
	}
	storageObjects := storageService.client.Bucket(storageService.bucketName).Objects(context, query)
	_, err := storageObjects.Next()
	if err == iterator.Done {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("gcloudstorage - FileExists: can't iterate over storage objects; %w", err)
	}
	return true, nil
}

func (storageService *gCloudStorageService) CreateBucket(context context.Context) error {
	bucket := storageService.client.Bucket(storageService.bucketName)
	if err := bucket.Create(context, storageService.projectName, nil); err != nil {
		return fmt.Errorf("gcloudstorage - CreateBucket: can't create bucket; %w", err)
	}
	return nil
}

func (storageService *gCloudStorageService) makeSignedUrlOptions(expiresInNumberOfMinutes int, insecure bool) *storage.SignedURLOptions {
	opts := &storage.SignedURLOptions{
		Method:   "GET",
		Expires:  time.Now().Add(time.Duration(expiresInNumberOfMinutes) * time.Minute),
		Scheme:   storage.SigningSchemeV4,
		Insecure: insecure,
	}
	if storageService.credentials != nil {
		opts.PrivateKey = storageService.credentials.PrivateKey
		opts.GoogleAccessID = storageService.credentials.AccessID
	}
	return opts
}
func (storageService *gCloudStorageService) GetExpiringUrl(tenant, fileName string, expiresInNumberOfMinutes int, insecure bool) (string, error) {
	return storageService.MakeSignedUrl(tenant, fileName, expiresInNumberOfMinutes, insecure)
}

func (storageService *gCloudStorageService) MakeSignedUrl(tenant, fileName string, expiresInNumberOfMinutes int, insecure bool) (string, error) {
	opts := storageService.makeSignedUrlOptions(expiresInNumberOfMinutes, insecure)
	queryParams := make(map[string][]string)
	fileExtension := filepath.Ext(fileName)
	queryParams["response-content-disposition"] = append(queryParams["response-content-disposition"], "attachment")
	queryParams["response-content-type"] = append(queryParams["response-content-disposition"], mime.TypeByExtension(fileExtension))
	opts.QueryParameters = queryParams
	signedUrl, err := storageService.client.Bucket(storageService.bucketName).SignedURL(fmt.Sprintf("%s/%s", tenant, fileName), opts)
	if err != nil {
		return "", fmt.Errorf("gcloudstorage - MakeSignedUrl: can't make signed url; %w", err)
	}
	return signedUrl, nil
}
