package files

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"

	"github.com/marcinlovescode/go-clean-fileupload/internal/pkg/logger"
)

type FileRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type File struct {
	ID   int       `json:"id" example:"1" format:"int64"`
	UUID uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

func AppendFileRoutes(handler *gin.RouterGroup, logger logger.Logger) {
	routerGroup := handler.Group("/files")
	routerGroup.GET("/ping", pingHandler)
	routerGroup.GET("/:id", showFile(logger))
	routerGroup.POST("/", uploadFile(logger))
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

// ShowFile
//
// @Summary     Show file
// @Description Get file details by id
// @ID          get-file-by-id
// @Tags  	    files
// @Accept      json
// @Produce     json
// @Param		id	path string	true "File ID"
// @Success     200 {object} File
// @Failure     400 {string} Error
// @Failure     500 {string} Error
// @Router      /files/{id} [get]
func showFile(logger logger.Logger) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var file FileRequest
		if err := ginCtx.ShouldBindUri(&file); err != nil {
			logger.Error(err, "http - files")
			ginCtx.String(http.StatusBadRequest, err.Error())
			return
		}
		parsedUUID, err := uuid.Parse(file.ID)
		if err != nil {
			logger.Error(err, "http - files")
			ginCtx.String(http.StatusBadRequest, err.Error())
			return
		}
		ginCtx.JSON(http.StatusOK, File{ID: 1, UUID: parsedUUID})
	}
}

// uploadFile godoc
//
// @Summary     Upload file
// @Description Upload file and attach it to the object of id
// @ID          upload-file
// @Tags  	    files
// @Accept      multipart/form-data
// @Produce     json
// @Param		file formData file true "File"
// @Param		referenceObjectId formData string true "Reference Object ID"
// @Success     200 {object} File
// @Failure     500 {string} Error
// @Router      /files/ [post]
func uploadFile(logger logger.Logger) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		referenceObjectId := ginCtx.PostForm("referenceObjectId")
		if referenceObjectId == "" {
			logger.Debug("http - files - uploadFile - referenceObjectId is empty")
			ginCtx.String(http.StatusBadRequest, "referenceObjectId is empty")
		}

		file, err := ginCtx.FormFile("file")
		if err != nil {
			logger.Debug(err, "http - files - uploadFile")
			ginCtx.String(http.StatusBadRequest, "corrupted file")
			return
		}

		filename := filepath.Base(file.Filename)
		if err := ginCtx.SaveUploadedFile(file, filename); err != nil {
			logger.Debug(err, "http - files - uploadFile")
			ginCtx.String(http.StatusBadRequest, "corrupted file")
			return
		}

		ginCtx.JSON(http.StatusOK, "Empty")
	}
}
