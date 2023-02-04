package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/config"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/app"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/models"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/gcloudstorage"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/logger"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/tests/doubles"
)

func createStubLogger() logger.Logger {
	return &doubles.VoidLogger{}
}

func TestPingRoute(t *testing.T) {
	// Arrange
	configPath := "./../config/config.yml"
	appConfig, err := config.NewConfig(configPath)
	if err != nil {
		t.FailNow()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	router, err := app.CreateHttpHandlers(ctx, createStubLogger(), appConfig)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/files/ping", nil)
	// Act
	router.ServeHTTP(w, req)
	// Assert
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Pong", w.Body.String())
}

func TestCanUploadFileAndDownloadFromStorageByUrl(t *testing.T) {
	// Arrange
	referenceObjectId := uuid.New().String()
	filename := fmt.Sprintf("textfile_%d_1.txt", uuid.New().ID())
	textContent := fmt.Sprintf("textcontent-%s", uuid.New().String())
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	appConfig := getConfig(t)
	setupGCSEmulator(t, ctx, appConfig)
	createRequest := makeUploadFileRequest(t, referenceObjectId, filename, textContent)
	getAttachmentsRequest := makeGetAttachmentsRequest(t, referenceObjectId)
	router, err := app.CreateHttpHandlers(ctx, createStubLogger(), appConfig)
	if err != nil {
		t.Fail()
	}
	createFileRequestRecorder := httptest.NewRecorder()
	getAttachmentsRequestRecorder := httptest.NewRecorder()
	// Act + Assert
	router.ServeHTTP(createFileRequestRecorder, createRequest)
	router.ServeHTTP(getAttachmentsRequestRecorder, getAttachmentsRequest)
	assert.Equal(t, http.StatusNoContent, createFileRequestRecorder.Code)
	assert.Equal(t, http.StatusOK, getAttachmentsRequestRecorder.Code)
	var attachments []models.Attachment
	if err := json.NewDecoder(getAttachmentsRequestRecorder.Body).Decode(&attachments); err != nil {
		t.FailNow()
	}
	attachment := attachments[0]
	getFileRequestResponse, err := http.Get(attachment.Url)
	if err != nil {
		t.FailNow()
	}
	responseBody, err := io.ReadAll(getFileRequestResponse.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = getFileRequestResponse.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, http.StatusOK, getFileRequestResponse.StatusCode)
	assert.Equal(t, textContent, string(responseBody))
}

func TestCanRemoveFile(t *testing.T) {
	// Arrange
	referenceObjectId := uuid.New().String()
	filename := fmt.Sprintf("textfile_%d_1.txt", uuid.New().ID())
	textContent := fmt.Sprintf("textcontent-%s", uuid.New().String())
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	appConfig := getConfig(t)
	setupGCSEmulator(t, ctx, appConfig)
	createRequest := makeUploadFileRequest(t, referenceObjectId, filename, textContent)
	getAttachmentsRequest := makeGetAttachmentsRequest(t, referenceObjectId)
	router, err := app.CreateHttpHandlers(ctx, createStubLogger(), appConfig)
	if err != nil {
		t.Fail()
	}
	createFileRequestRecorder := httptest.NewRecorder()
	getAttachmentsRequestRecorder := httptest.NewRecorder()
	removeFileRequestRecorder := httptest.NewRecorder()

	// Act + Assert
	router.ServeHTTP(createFileRequestRecorder, createRequest)
	router.ServeHTTP(getAttachmentsRequestRecorder, getAttachmentsRequest)
	assert.Equal(t, http.StatusNoContent, createFileRequestRecorder.Code)
	assert.Equal(t, http.StatusOK, getAttachmentsRequestRecorder.Code)
	var attachments []models.Attachment
	if err := json.NewDecoder(getAttachmentsRequestRecorder.Body).Decode(&attachments); err != nil {
		t.FailNow()
	}
	attachment := attachments[0]
	getFileRequestResponse := makeRemoveFileRequest(t, attachment.ID)
	router.ServeHTTP(removeFileRequestRecorder, getFileRequestResponse)
	assert.Equal(t, http.StatusNoContent, removeFileRequestRecorder.Code)
}

func setupBucket(ctx context.Context, projectName, bucketName string) error {
	gcloudClient, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	fileService := gcloudstorage.NewGCloudStorageService(gcloudClient, projectName, bucketName, nil)
	return fileService.CreateBucket(ctx)
}

func getConfig(t *testing.T) *config.Config {
	configPath := "./../config/config.yml"
	appConfig, err := config.NewConfig(configPath)
	if err != nil {
		t.FailNow()
	}
	appConfig.GCloudStorage.BucketName = fmt.Sprintf("files-%d-1", uuid.New().ID())
	return appConfig
}

func setupGCSEmulator(t *testing.T, ctx context.Context, config *config.Config) {
	err := os.Setenv("STORAGE_EMULATOR_HOST", fmt.Sprintf("http://localhost:%d", config.GCloudStorage.EmulatorPort))
	if err != nil {
		t.FailNow()
	}
	err = setupBucket(ctx, config.GCloudStorage.ProjectName, config.GCloudStorage.BucketName)
	if err != nil {
		t.FailNow()
	}
}

func makeUploadFileRequest(t *testing.T, referenceObjectId, filename, textContent string) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.WriteField("referenceObjectId", referenceObjectId)
	if err != nil {
		t.FailNow()
	}
	part, _ := writer.CreateFormFile("file", filename)
	_, err = part.Write([]byte(textContent))
	if err != nil {
		t.FailNow()
	}
	err = writer.Close()
	if err != nil {
		t.FailNow()
	}
	req, err := http.NewRequest("POST", "/api/files", body)
	if err != nil {
		t.FailNow()
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func makeGetAttachmentsRequest(t *testing.T, referenceObjectId string) *http.Request {
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/files/reference/%s", referenceObjectId), nil)
	if err != nil {
		t.FailNow()
	}
	return req
}

func makeRemoveFileRequest(t *testing.T, fileId string) *http.Request {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/files/%s", fileId), nil)
	if err != nil {
		t.FailNow()
	}
	return req
}
