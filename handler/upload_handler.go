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

// Upload godoc
//
//	@Summary		Upload file(s)
//	@Description	Upload single or multiple images using multipart/form-data (key: files)
//	@Tags			Uploads
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			files	formData	file	true	"Upload file(s)"
//	@Success		200		{object}	dto.UploadSuccessSingleResponse
//	@Success		200		{object}	dto.UploadSuccessMultipleResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/uploads [post]
func (uh *uploadHandler) Upload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil || form.File == nil {
		uh.logger.Error("No files found in multipart form", zap.Error(err))
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_NO_FILES_UPLOADED, "no file(s) uploaded")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["files"]

	// upload with key "files"
	if len(files) == 0 {
		file, err := ctx.FormFile("files")
		if err != nil {
			uh.logger.Error("Failed to get file from form", zap.Error(err))
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_NO_FILES_UPLOADED, "no file(s) uploaded")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		files = []*multipart.FileHeader{file}
	}

	// call service
	uploadedImages, err := uh.uploadService.Upload(ctx, files)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_FILES, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// kalau hanya 1 file → balikin object image
	if len(uploadedImages) == 1 {
		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILE, uploadedImages[0])
		ctx.JSON(http.StatusOK, res)
		return
	}

	// kalau banyak file → balikin array object image
	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILES, uploadedImages)
	ctx.JSON(http.StatusOK, res)
}
