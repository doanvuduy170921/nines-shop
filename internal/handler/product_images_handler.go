package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
	"os"
	"path"
	"strconv"
)

type ProductImagesHandler struct {
	service service.ProductImagesService
}

func NewProductImagesHandler(service service.ProductImagesService) *ProductImagesHandler {
	return &ProductImagesHandler{
		service: service,
	}
}

func (ph *ProductImagesHandler) UploadImage(ctx *gin.Context) {
	imgage, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(400, validation.HandlerValidationError(err))
		return
	}
	if imgage.Size > 5<<20 {
		utils.ResponseError(ctx, errors.New("file size too big"))
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)

	dst := fmt.Sprintf("./uploads/%s", path.Base(imgage.Filename))

	err = ctx.SaveUploadedFile(imgage, dst)
	if err != nil {
		ctx.JSON(400, validation.HandlerValidationError(err))
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, imgage.Filename, "Upload File Success")

}

func (ph *ProductImagesHandler) MultipleUploadImages(ctx *gin.Context) {
	idStr := ctx.Param("product_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result, err := ph.service.SaveUploadFile(ctx, id)
	if err != nil && result == nil {
		// Lỗi nghiêm trọng (không parse được form, etc)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := http.StatusOK
	if result.SuccessCount == 0 {
		status = http.StatusBadRequest
	} else if result.FailedCount > 0 {
		status = http.StatusPartialContent // 206
	}

	ctx.JSON(status, gin.H{
		"message":       "Upload completed",
		"saved_images":  result.SavedImages,
		"failed_files":  result.FailedFiles,
		"success_count": result.SuccessCount,
		"failed_count":  result.FailedCount,
	})
}
