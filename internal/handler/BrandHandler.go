package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
)

type BrandHandler struct {
	service service.BrandService
}

func NewBrandHandler(service service.BrandService) *BrandHandler {
	return &BrandHandler{
		service: service,
	}
}

func (ch *BrandHandler) GetAllBrand(ctx *gin.Context) {
	categories, err := ch.service.GetAll(ctx)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, categories, "Get Brand Success")
}
