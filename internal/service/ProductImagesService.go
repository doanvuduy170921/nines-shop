package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
)

type productImagesService struct {
	repo repository.ProductImagesRepository
	pu   ProductUpdater
}

type ProductUpdater interface {
	GetProductById(ctx context.Context, id int32) (sqlc.Product, error)
	UpdateThumbnail(ctx context.Context, thumbnail string, id int32) (sqlc.Product, error)
}

func NewProductImagesService(repo repository.ProductImagesRepository, pu ProductUpdater) ProductImagesService {
	return &productImagesService{
		repo: repo,
		pu:   pu,
	}

}

func (ps *productImagesService) SaveUploadFile(ctx *gin.Context, id int64) (*dto.UploadResult, error) {
	context := ctx.Request.Context()

	const publicURL = "http://localhost:8080/api/v1/images/"

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	images := form.File["images"]
	if len(images) == 0 {
		return nil, errors.New("no file provided")
	}

	result := &dto.UploadResult{
		SavedImages: make([]sqlc.ProductImage, 0),
		FailedFiles: make([]dto.UploadError, 0),
	}

	for _, image := range images {
		filename, err := utils.ValidateAndSaveFile(image, "./uploads")
		if err != nil {
			result.FailedFiles = append(result.FailedFiles, dto.UploadError{
				Filename: image.Filename,
				Error:    err.Error(),
			})
			result.FailedCount++
			continue
		}

		publicImageURL := publicURL + filename

		savedImage, err := ps.repo.Save(context, sqlc.SaveAndUploadImgParams{
			ProductID: id,
			ImageUrl:  publicImageURL,
		})
		if err != nil {
			result.FailedFiles = append(result.FailedFiles, dto.UploadError{
				Filename: image.Filename,
				Error:    "DB error: " + err.Error(),
			})
			result.FailedCount++
			continue
		}

		result.SavedImages = append(result.SavedImages, savedImage)
		result.SuccessCount++
	}

	thumbnailProduct := result.SavedImages[0].ImageUrl // lấy ảnh đầu tiên của product set cho thumbnail
	_, err = ps.pu.UpdateThumbnail(ctx, thumbnailProduct, int32(id))
	if err != nil {
		log.Println("Error updating thumbnail: " + err.Error())
		return result, err
	}
	if result.SuccessCount == 0 {
		return result, errors.New("all uploads failed")
	}

	return result, nil
}
