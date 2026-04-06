package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
)

type MediaHandler struct {
	mediaService service.MediaService
}

func NewMediaHandler(ms service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: ms}
}

func (h *MediaHandler) UploadImages(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn ít nhất 1 ảnh"})
		return
	}

	urls, errs := h.mediaService.UploadMultiple(files)

	// Trả về kết quả cho Frontend
	if len(errs) > 0 {
		var msgErr string
		for _, err := range errs {
			msgErr += err.Error() + ";"
		}
		utils.ResponseError(ctx, utils.NewError(http.StatusBadRequest, "Upload loi :"+msgErr))
	}
	utils.ResponseSuccess(ctx, http.StatusOK, urls, "Upload images successfully")
}
