package files

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/files/models"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/pkg/logger"
)

type FileRequest struct {
	ID string `uri:"id" binding:"required"`
}

type Attachments []models.Attachment

func AppendFileRoutes(handler *gin.RouterGroup, logger logger.Logger, useCase UseCase) {
	routerGroup := handler.Group("/files")
	routerGroup.GET("/ping", pingHandler)
	routerGroup.GET("/reference/:id", showFiles(logger, useCase))
	routerGroup.POST("", uploadFile(logger, useCase))
	routerGroup.DELETE(":id", deleteFile(logger, useCase))

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

// showFiles godoc
//
// @Summary     Show files
// @Description Get file by reference id
// @ID          get-file-by-reference-id
// @Tags  	    files
// @Accept      json
// @Produce     json
// @Param		id	path string	true "FileDto ID"
// @Success     200 {object} Attachments
// @Failure     400 {string} Error
// @Failure     500 {string} Error
// @Router      /files/reference/{id} [get]
func showFiles(logger logger.Logger, useCase UseCase) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Copy()
		var file FileRequest
		if err := ginCtx.ShouldBindUri(&file); err != nil {
			logger.Debug(err, "files - showFiles")
			ginCtx.String(http.StatusBadRequest, err.Error())
			return
		}
		result, err := useCase.ListBy(ctx, "tenant1", file.ID)
		if err != nil {
			logger.Error(err, "files - showFiles")
			ginCtx.String(http.StatusInternalServerError, "can't get files")
			return
		}
		ginCtx.JSON(http.StatusOK, *result)
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
// @Param		file formData file true "FileDto"
// @Param		referenceObjectId formData string true "Reference Object ID"
// @Success     204
// @Failure     500 {string} Error
// @Router      /files/ [post]
func uploadFile(logger logger.Logger, useCase UseCase) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Copy()
		referenceObjectId := ginCtx.PostForm("referenceObjectId")
		if referenceObjectId == "" {
			logger.Debug("files - uploadFile - referenceObjectId is empty")
			ginCtx.String(http.StatusBadRequest, "referenceObjectId is empty")
		}

		file, err := ginCtx.FormFile("file")
		if err != nil {
			logger.Debug(err, "files - uploadFile")
			ginCtx.String(http.StatusBadRequest, "corrupted file")
			return
		}
		filename := filepath.Base(file.Filename)
		fileHandler, err := file.Open()
		if err != nil {
			logger.Debug(err, "files - uploadFile")
			ginCtx.String(http.StatusBadRequest, "corrupted file")
			return
		}
		err = useCase.UploadFile(ctx, "tenant1", models.UploadFileCommand{
			CreatorId:   "UserId",
			FileName:    filename,
			ReferenceID: referenceObjectId,
			File:        fileHandler,
		})
		if err != nil {
			logger.Error(err, "files - uploadFile")
			ginCtx.String(http.StatusBadRequest, "can't upload file")
			return
		}
		ginCtx.Status(http.StatusNoContent)
	}
}

// deleteFile godoc
//
// @Summary     Remove file
// @Description Remove file by id
// @ID          remove-file-by-id
// @Tags  	    files
// @Accept      json
// @Produce     json
// @Param		id	path string	true "FileID"
// @Success     204
// @Failure     400 {string} Error
// @Failure     500 {string} Error
// @Router      /files/{id} [delete]
func deleteFile(logger logger.Logger, useCase UseCase) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Copy()
		var file FileRequest
		if err := ginCtx.ShouldBindUri(&file); err != nil {
			logger.Debug(err, "files - deleteFile")
			ginCtx.String(http.StatusBadRequest, err.Error())
			return
		}
		err := useCase.DeleteFile(ctx, "tenant1", file.ID)
		if err != nil {
			logger.Error(err, "files - deleteFile")
			ginCtx.String(http.StatusInternalServerError, "can't remove file")
			return
		}
		ginCtx.Status(http.StatusNoContent)
	}
}
