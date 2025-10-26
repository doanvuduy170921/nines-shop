package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var input dto.CreateCategoryParamDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, validation.HandlerValidationError(err))
		return
	}

	category, err := ch.service.Create(ctx, input.Name)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, category, "Create category Success")
}

func (ch *CategoryHandler) GetAllCategory(ctx *gin.Context) {
	categories, err := ch.service.GetAll(ctx)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, categories, "Get Category Success")
}
