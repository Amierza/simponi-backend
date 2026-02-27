package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	IUploadHandler interface {
		Upload(ctx *gin.Context)
	}

	uploadHandler struct {
		uploadService service.IUploadService
		logger        *zap.Logger
	}
)

func NewUploadHandler(uploadService service.IUploadService, logger *zap.Logger) *uploadHandler {
	return &uploadHandler{
		uploadService: uploadService,
		logger:        logger,
	}
}

func (uh *uploadHandler) Upload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil || form.File == nil {
		uh.logger.Error("No files found in multipart form", zap.Error(err))
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_NO_FILES_UPLOADED, "no file(s) uploaded", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["files"]

	// upload with key "files"
	if len(files) == 0 {
		file, err := ctx.FormFile("files")
		if err != nil {
			uh.logger.Error("Failed to get file from form", zap.Error(err))
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_NO_FILES_UPLOADED, "no file(s) uploaded", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		files = []*multipart.FileHeader{file}
	}

	// call service
	uploadedURLs, err := uh.uploadService.Upload(ctx, files)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_FILES, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// kalau hanya 1 file → balikin string saja
	if len(uploadedURLs) == 1 {
		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILE, uploadedURLs[0])
		ctx.JSON(http.StatusOK, res)
		return
	}

	// kalau banyak file → balikin array
	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILES, uploadedURLs)
	ctx.JSON(http.StatusOK, res)
}
