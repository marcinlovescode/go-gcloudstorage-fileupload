package http

import (
	"github.com/gin-gonic/gin"

	_ "github.com/marcinlovescode/go-gcloudstorage-fileupload/docs"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/logger"
)

// NewGinHttpRouter -.
// Swagger spec:
// @title       FileRequest Upload Service
// @description Filed Upload Service that stores files associated by reference id
// @version     1.0
// @host        localhost:8080
// @BasePath    /api
func NewGinHttpRouter(logger logger.Logger, useCase files.UseCase, handler *gin.Engine) error {
	appendApiRoutes(handler, logger, useCase)
	return nil
}

func appendApiRoutes(handler *gin.Engine, logger logger.Logger, filesUseCase files.UseCase) {
	routerGroup := handler.Group("/api")
	files.AppendFileRoutes(routerGroup, logger, filesUseCase)
}
