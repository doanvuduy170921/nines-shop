package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
)

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo}

}

func (ps *productService) CreateProduct(ctx *gin.Context, arg dto.CreateProductParamDto) (sqlc.Product, error) {
	context := ctx.Request.Context()

	sku := utils.GenProductSku(arg.Name)
	slug := utils.GenProductSlug(arg.Name)
	productDto := dto.MapProductDtoToParams(arg)

	product, err := ps.repo.CreateProduct(context, sqlc.CreateProductParams{
		Sku:              sku,
		Slug:             slug,
		Name:             productDto.Name,
		BrandID:          productDto.BrandID,
		CategoryID:       productDto.CategoryID,
		Description:      productDto.Description,
		ShortDescription: productDto.ShortDescription,
		Price:            productDto.Price,
		DiscountPrice:    productDto.DiscountPrice,
		StockQuantity:    productDto.StockQuantity,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.Product{}, utils.WrapError(err, "Product is exists", utils.ErrorCodeBadRequest)
		}

		return sqlc.Product{}, err
	}
	return product, err
}

func (ps *productService) GetAllByFilter(ctx *gin.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) ([]sqlc.Product, int64, error) {
	context := ctx.Request.Context()
	products, err := ps.repo.GetAllProductByFilter(context, limit, page, categoryId, minPrice, maxPrice, search, status)
	if err != nil {
		return nil, 0, err
	}
	count, err := ps.repo.CountProduct(context, limit, page, categoryId, minPrice, maxPrice, search, status)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}
