package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/pagination"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
	"strconv"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (ph *ProductHandler) GetAllByFilter(ctx *gin.Context) {

	listProducts, err := ph.service.GetAllProductByFilter(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, listProducts, "Get All Product By Filter successfully")

}

func (ph *ProductHandler) AddProduct(ctx *gin.Context) {
	var input dto.AddProductRequestDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, validation.HandlerValidationError(err))
		return
	}
	product, err := ph.service.AddProduct(ctx.Request.Context(), input)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, product, "Add Product Success")
}

func (ph *ProductHandler) GetListVariantByPid(ctx *gin.Context) {
	pIdStr := ctx.Param("product_id")
	pId, err := strconv.Atoi(pIdStr)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	listVariants, err := ph.service.GetListVariantByPid(ctx.Request.Context(), int64(pId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, listVariants, "Get List Variant By Pid Success")
}

func (ph *ProductHandler) GetTop3Thumbnail(ctx *gin.Context) {
	data, err := ph.service.GetTop3Thumbnail(ctx.Request.Context())
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, data, "Get Data for page home Success")

}

func (ph *ProductHandler) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	product, err := ph.service.GetProductBySlug(ctx.Request.Context(), slug)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, product, "Get Product By Slug Success")
}

func (ph *ProductHandler) GetListProducts(ctx *gin.Context) {
	var arg dto.GetProductsRequest
	if err := ctx.ShouldBindQuery(&arg); err != nil {
		ctx.JSON(400, validation.HandlerValidationError(err))
		return
	}

	products, err := ph.service.GetListProducts(ctx.Request.Context(), arg)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	res := pagination.PaginationRes{
		Data: products,
	}

	utils.ResponseSuccess(ctx, http.StatusOK, res, "Get Products Success")

}
