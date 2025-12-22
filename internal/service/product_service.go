package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
)

type productService struct {
	repo repository.ProductRepository
	iu   ImagesUpdater
}

type ImagesUpdater interface {
	GetImagesByProductId(ctx context.Context, id int32) ([]string, error)
}

func NewProductService(repo repository.ProductRepository, iu ImagesUpdater) ProductService {
	return &productService{
		repo: repo,
		iu:   iu,
	}

}

func (ps *productService) CreateProduct(ctx *gin.Context, arg dto.CreateProductParamDto) (sqlc.Product, error) {
	c := ctx.Request.Context()

	sku := utils.GenProductSku(arg.Name)
	slug := utils.GenProductSlug(arg.Name)
	productDto := dto.MapProductDtoToParams(arg)

	product, err := ps.repo.CreateProduct(c, sqlc.CreateProductParams{
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

func (ps *productService) GetAllByFilter(ctx *gin.Context, limit, page, categoryId, minPrice, maxPrice, brandId int32, search, status string) ([]sqlc.GetAllProductByFilterRow, int64, error) {
	c := ctx.Request.Context()
	products, err := ps.repo.GetAllProductByFilter(c, limit, page, categoryId, minPrice, maxPrice, brandId, search, status)
	if err != nil {
		return nil, 0, err
	}
	count, err := ps.repo.CountProduct(c, limit, page, categoryId, minPrice, maxPrice, search, status)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (ps *productService) GetProductByCategoryId(ctx *gin.Context, id int32) ([]sqlc.GetProductByCategoryIdRow, error) {
	c := ctx.Request.Context()
	products, err := ps.repo.GetProductByCategoryId(c, id)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *productService) GetProductBySlug(ctx *gin.Context, slug string) (sqlc.GetProductBySlugRow, error) {
	c := ctx.Request.Context()
	product, err := ps.repo.GetProductBySlug(c, slug)
	if err != nil {
		return sqlc.GetProductBySlugRow{}, err
	}
	return product, nil
}

func (ps *productService) GetImagesBySlug(c *gin.Context, slug string) ([]string, error) {
	ctx := c.Request.Context()
	product, err := ps.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	images, err := ps.iu.GetImagesByProductId(ctx, int32(product.ID))
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (ps *productService) GetTop8ProductSeller(c *gin.Context, cateID *int32) ([]sqlc.GetTop8ProductSellerRow, error) {
	ctx := c.Request.Context()
	products, err := ps.repo.GetTop8ProductSeller(ctx, cateID)
	if err != nil {
		return nil, utils.WrapError(err, "Get Top 3 seller product fail", http.StatusBadRequest)
	}
	return products, nil
}
