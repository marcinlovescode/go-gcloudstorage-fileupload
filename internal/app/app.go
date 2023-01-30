// Package app configures and runs application.
package app

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/marcinlovescode/go-clean-fileupload/config"
	httpRouter "github.com/marcinlovescode/go-clean-fileupload/internal/http"
	"github.com/marcinlovescode/go-clean-fileupload/internal/pkg/httpserver"
	"github.com/marcinlovescode/go-clean-fileupload/internal/pkg/logger"
)

func Run(cfg *config.Config) {
	zerologLogger := logger.NewZerologLogger(cfg.Log.Level)
	handler := CreateHttpHandlers(zerologLogger)
	AppendSwagger(handler)
	httpRouter.NewGinHttpRouter(handler, zerologLogger)
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

	err := httpServer.Shutdown()
	if err != nil {
		zerologLogger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func CreateHttpHandlers(logger logger.Logger) *gin.Engine {
	handler := gin.Default()
	httpRouter.NewGinHttpRouter(handler, logger)
	return handler
}

func AppendSwagger(handler *gin.Engine) {
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
}
