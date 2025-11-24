package handler

import (
	"github.com/gin-gonic/gin"
	"math"
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

func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	var input dto.CreateProductParamDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, validation.HandlerValidationError(err))
		return
	}

	product, err := ph.service.CreateProduct(ctx, input)
	if err != nil {

		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, product, "Create Product Success")
}

func (ph *ProductHandler) GetAllByFilter(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	pageStr := ctx.DefaultQuery("page", "1")
	searchStr := ctx.Query("search")
	statusStr := ctx.Query("status")
	minPriceStr := ctx.DefaultQuery("min_price", "0")
	maxPriceStr := ctx.DefaultQuery("max_price", "999999999")
	categoryIdStr := ctx.Query("category_id")

	var categoryId int32
	if categoryIdStr != "" {
		catId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			utils.ResponseError(ctx, err)
			return
		}
		categoryId = int32(catId)
	} else {
		// Nếu không truyền categoryId, set về 0 để query lấy tất cả
		categoryId = 0
	}

	limit := validation.ConvertAndValidateIntParam(ctx, limitStr, 10)
	page := validation.ConvertAndValidateIntParam(ctx, pageStr, 1)
	minPrice := validation.ConvertAndValidateIntParam(ctx, minPriceStr, 0)
	maxPrice := validation.ConvertAndValidateIntParam(ctx, maxPriceStr, 999999999)

	products, count, err := ph.service.GetAllByFilter(ctx, limit, page, categoryId, minPrice, maxPrice, searchStr, statusStr)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	productsRes := dto.MapSqlcProductsToResponse(products)

	paginationRes := pagination.PaginationRes{
		Data:      productsRes,
		Total:     int32(count),
		Page:      int32(page),
		Limit:     int32(limit),
		TotalPage: int32(math.Ceil(float64(count) / float64(limit))),
	}

	utils.ResponseSuccess(ctx, http.StatusOK, paginationRes, "Get All Product Success")

}

func (ph *ProductHandler) GetAllByCateId(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	products, err := ph.service.GetProductByCategoryId(ctx, int32(id))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	productRes := dto.MapProductsCateIdToRes(products)
	utils.ResponseSuccess(ctx, http.StatusOK, productRes, "Get All Product Success")
}

func (ph *ProductHandler) GetAllBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	product, err := ph.service.GetProductBySlug(ctx, slug)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	productRes := dto.MapProductSlugToRes(product)
	utils.ResponseSuccess(ctx, http.StatusOK, productRes, "Get All Product Success")
}

func (ph *ProductHandler) GetImagesBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	images, err := ph.service.GetImagesBySlug(ctx, slug)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, images, "Get Images Success")
}
