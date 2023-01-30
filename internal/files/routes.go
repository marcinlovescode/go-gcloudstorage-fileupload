package files

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/marcinlovescode/go-clean-fileupload/internal/pkg/logger"
)

func AppendFileRoutes(handler *gin.RouterGroup, logger logger.Logger) {
	routerGroup := handler.Group("/files")
	routerGroup.GET("/ping", pingHandler)
}

// pingHandler godoc
//
// @Summary     Ping-pong healthcheck
// @Description Returns Pong when service is alive
// @ID          files-ping-pong
// @Tags  	    files
// @Accept      json
// @Produce     json
// @Success     200 {string} Pong
// @Failure     500 {string} Error
// @Router      /files/ping [get]
func pingHandler(ginCtx *gin.Context) {
	ginCtx.String(http.StatusOK, "Pong")
}
