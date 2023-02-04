package app

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/config"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/adapters"
	httpRouter "github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/http"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/gcloudstorage"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/httpserver"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/logger"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	zerologLogger := logger.NewZerologLogger(cfg.Log.Level)
	handler, err := CreateHttpHandlers(ctx, zerologLogger, cfg)
	if err != nil {
		zerologLogger.Error(fmt.Errorf("app - Run - init: %w", err))
		cancel()
		return
	}
	appendSwagger(handler)
	httpServer := httpserver.New(handler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ShutdownTimeout(time.Duration(cfg.HTTP.ShutdownTimeout)),
		httpserver.WriteTimeout(time.Duration(cfg.HTTP.WriteTimeout)),
		httpserver.ReadTimeout(time.Duration(cfg.HTTP.ReadTimeout)))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		zerologLogger.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		zerologLogger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	cancel()
	err = httpServer.Shutdown()
	if err != nil {
		zerologLogger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func CreateHttpHandlers(ctx context.Context, logger logger.Logger, cfg *config.Config) (*gin.Engine, error) {
	handler := gin.Default()
	useCase, err := createFilesUseCase(ctx, cfg.GCloudStorage)
	if err != nil {
		return nil, fmt.Errorf("app - createHttpHandlers: can't create Files UseCase; %w", err)
	}
	err = httpRouter.NewGinHttpRouter(logger, useCase, handler)
	if err != nil {
		return nil, fmt.Errorf("app - createHttpHandlers: can't create router; %w", err)
	}
	return handler, nil
}

func appendSwagger(handler *gin.Engine) {
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
}

func createFilesUseCase(ctx context.Context, gcloudConfig config.GCloudStorage) (files.UseCase, error) {
	gcloudClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("http - router - newFilesUseCase: can't create Google Cloud Storage Client; %w", err)
	}
	var gCSCredentials *gcloudstorage.GCloudAuthCredentials
	if gcloudConfig.UseCredentials {
		privateKey, err := b64.StdEncoding.DecodeString(gcloudConfig.PrivateKeyBase64)
		if err != nil {
			return nil, fmt.Errorf("http - router - newFilesUseCase: can't decode private key; %w", err)
		}
		gCSCredentials = &gcloudstorage.GCloudAuthCredentials{
			AccessID:   gcloudConfig.AccessId,
			PrivateKey: privateKey,
		}
	}
	fileService := gcloudstorage.NewGCloudStorageService(gcloudClient, gcloudConfig.ProjectName, gcloudConfig.BucketName, gCSCredentials)
	fileRepository, err := adapters.NewInMemoryFilesRepository(ctx)
	idGen := adapters.NewGuidBasedIdGenerator()

	return files.NewDefaultFilesUseCase(fileService, fileRepository, idGen, gcloudConfig)
}
