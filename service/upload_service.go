package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IUploadService interface {
		// public function
		Upload(ctx context.Context, files []*multipart.FileHeader) ([]dto.UploadImageResponse, error)
		// private / helper function
		saveUploadedFile(file *multipart.FileHeader, savePath string) error
		createFile(path string) (*os.File, error)
		copyFile(dst *os.File, src multipart.File) (int64, error)
	}

	uploadService struct {
		productRepo repository.IProductRepository
		logger      *zap.Logger
	}
)

func NewUploadService(productRepo repository.IProductRepository, logger *zap.Logger) *uploadService {
	return &uploadService{
		productRepo: productRepo,
		logger:      logger,
	}
}

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
}

// Upload bisa handle single atau multiple file
func (us *uploadService) Upload(ctx context.Context, files []*multipart.FileHeader) ([]dto.UploadImageResponse, error) {
	if len(files) == 0 {
		us.logger.Warn("Upload attempted with no files")
		return nil, dto.ErrNoFilesUploaded
	}

	var uploadedImages []dto.UploadImageResponse
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExt[ext] {
			us.logger.Warn("Invalid file type",
				zap.String("filename", file.Filename),
				zap.String("extension", ext),
			)
			return nil, dto.ErrInvalidFileType
		}

		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		storagePath := filepath.Join("uploads", newFileName)

		// Simpan file (local)
		if err := us.saveUploadedFile(file, storagePath); err != nil {
			us.logger.Error("Failed to save uploaded file",
				zap.String("filename", file.Filename),
				zap.String("path", storagePath),
				zap.Error(err),
			)
			return nil, dto.ErrSaveFile
		}

		newImage, err := us.productRepo.CreateProductImage(ctx, nil, &entity.ProductImage{
			ID:       uuid.New(),
			ImageURL: storagePath,
		})
		if err != nil {
			us.logger.Error("Failed to store uploaded image metadata", zap.String("path", storagePath), zap.Error(err))
			return nil, dto.ErrInternal
		}

		uploadedImages = append(uploadedImages, dto.UploadImageResponse{
			ImageID:  newImage.ID,
			ImageURL: newImage.ImageURL,
		})
	}

	return uploadedImages, nil
}

func (us *uploadService) saveUploadedFile(file *multipart.FileHeader, savePath string) error {
	src, err := file.Open()
	if err != nil {
		us.logger.Error("Failed to open uploaded file",
			zap.String("filename", file.Filename),
			zap.Error(err),
		)
		return err
	}
	defer src.Close()

	dst, err := us.createFile(savePath)
	if err != nil {
		us.logger.Error("Failed to create destination file",
			zap.String("path", savePath),
			zap.Error(err),
		)
		return err
	}
	defer dst.Close()

	_, err = us.copyFile(dst, src)
	if err != nil {
		us.logger.Error("Error while copying file content",
			zap.String("destination", savePath),
			zap.Error(err),
		)
	}
	return err
}

// ini nanti kita ganti kalau mau langsung ke S3
func (us *uploadService) createFile(path string) (*os.File, error) {
	return os.Create(path)
}

func (us *uploadService) copyFile(dst *os.File, src multipart.File) (int64, error) {
	return io.Copy(dst, src)
}
