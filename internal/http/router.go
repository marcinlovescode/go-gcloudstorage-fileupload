package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/marcinlovescode/go-clean-fileupload/docs"
	"github.com/marcinlovescode/go-clean-fileupload/internal/files"
	"github.com/marcinlovescode/go-clean-fileupload/internal/pkg/logger"
)

// NewGinHttpRouter -.
// Swagger spec:
// @title       FileRequest Upload Service
// @description Filed Upload Service that stores files associated by reference id
// @version     1.0
// @host        localhost:8080
// @BasePath    /api
func NewGinHttpRouter(handler *gin.Engine, logger logger.Logger) {
	appendApiRoutes(handler, logger)
}

func appendApiRoutes(handler *gin.Engine, logger logger.Logger) {
	routerGroup := handler.Group("/api")
	files.AppendFileRoutes(routerGroup, logger)
}
