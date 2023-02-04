package tests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/gcloudstorage"
)

func TestCanUploadFileAndDownloadItBySignedUrl(t *testing.T) {
	// Arrange
	err := setupEmulator()
	requireNotError(t, err)
	projectName := "my-project"
	bucketName := fmt.Sprintf("bucket-%d-1", uuid.New().ID())
	fileName := "file.txt"
	tenantName := "tenant"
	fileContent := []byte("Hello!")
	pkey, err := os.ReadFile("./data/key")
	requireNotError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := createGCloudStorageClient(t, ctx)
	service := gcloudstorage.NewGCloudStorageService(client, projectName, bucketName, &gcloudstorage.GCloudAuthCredentials{
		AccessID:   "xxx@developer.gserviceaccount.com",
		PrivateKey: pkey,
	})

	// Act
	err = service.CreateBucket(ctx)
	requireNotError(t, err)
	err = service.UploadFile(ctx, tenantName, fileName, fileContent)
	requireNotError(t, err)
	url, err := service.MakeSignedUrl(tenantName, fileName, 15, true)
	requireNotError(t, err)

	// Assert
	resp, err := http.Get(url)
	requireNotError(t, err)
	bodyBytes, _ := io.ReadAll(resp.Body)
	require.Equal(t, fileContent, bodyBytes)
	err = resp.Body.Close()
	requireNotError(t, err)
}

func TestCanRemoveFileFromStorage(t *testing.T) {
	// Arrange
	err := setupEmulator()
	requireNotError(t, err)
	projectName := "my-project"
	bucketName := fmt.Sprintf("bucket-%d-2", uuid.New().ID())
	fileName := "file.txt"
	tenantName := "tenant"
	fileContent := []byte("Hello!")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := createGCloudStorageClient(t, ctx)
	service := gcloudstorage.NewGCloudStorageService(client, projectName, bucketName, nil)
	// Act
	err = service.CreateBucket(ctx)
	requireNotError(t, err)
	err = service.UploadFile(ctx, tenantName, fileName, fileContent)
	requireNotError(t, err)
	err = service.DeleteFile(ctx, tenantName, fileName)
	requireNotError(t, err)
	// Assert
	resp, err := service.FileExists(ctx, tenantName, fileName)
	requireNotError(t, err)
	require.Equal(t, false, resp)
}

func requireNotError(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func setupEmulator() error {
	return os.Setenv("STORAGE_EMULATOR_HOST", "http://localhost:9023")
}

func createGCloudStorageClient(t *testing.T, ctx context.Context) (client *storage.Client) {
	client, err := storage.NewClient(ctx)
	requireNotError(t, err)
	return client
}
