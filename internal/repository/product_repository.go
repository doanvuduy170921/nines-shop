package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type productRepository struct {
	DB *sqlc.Queries
}

func NewProductRepository(DB *sqlc.Queries) ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (pr *productRepository) CountProduct(ctx context.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) (int64, error) {
	count, err := pr.DB.CountProduct(ctx, sqlc.CountProductParams{
		Search:     search,
		Status:     status,
		CategoryID: categoryId,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *productRepository) GetProductById(ctx context.Context, id int32) (sqlc.Product, error) {
	return pr.DB.GetProductById(ctx, id)
}

func (pr *productRepository) UpdateThumbnail(ctx context.Context, thumbnail string, id int32) (sqlc.Product, error) {
	product, err := pr.DB.UpdateThumbnail(ctx, sqlc.UpdateThumbnailParams{
		Thumbnail: thumbnail,
		ID:        id,
	})
	if err != nil {
		return sqlc.Product{}, err
	}
	return product, nil
}

func (pr *productRepository) AddProduct(ctx context.Context, arg sqlc.AddProductParams) (sqlc.Product, error) {
	return pr.DB.AddProduct(ctx, arg)
}

func (pr *productRepository) AddProductSpec(ctx context.Context, arg sqlc.AddProductSpecParams) (sqlc.ProductSpecification, error) {
	return pr.DB.AddProductSpec(ctx, arg)
}

func (pr *productRepository) AddProductVariant(ctx context.Context, arg sqlc.AddProductVariantParams) (sqlc.ProductVariant, error) {
	return pr.DB.AddProductVariant(ctx, arg)
}

func (pr *productRepository) AddAttributesConf(ctx context.Context, arg sqlc.AddAttributesConfParams) (sqlc.AttributeConfig, error) {
	return pr.DB.AddAttributesConf(ctx, arg)
}

func (pr *productRepository) GetAllProductByFilter(ctx context.Context) ([]sqlc.GetAllProductByFilterRow, error) {
	return pr.DB.GetAllProductByFilter(ctx)
}

func (pr *productRepository) GetListVariantByPid(ctx context.Context, productID int64) ([]sqlc.GetListVariantByPidRow, error) {
	return pr.DB.GetListVariantByPid(ctx, productID)
}

func (pr *productRepository) GetTop3Thumbnail(ctx context.Context) ([]string, error) {
	return pr.DB.GetTop3Thumbnail(ctx)
}

func (pr *productRepository) GetTop3Trending(ctx context.Context, cateID *int32) ([]sqlc.GetTop3TrendingRow, error) {
	return pr.DB.GetTop3Trending(ctx, cateID)
}

func (pr *productRepository) GetProductBySlug(ctx context.Context, slug string) (sqlc.GetProductBySlugRow, error) {
	return pr.DB.GetProductBySlug(ctx, slug)
}

func (pr *productRepository) GetVariantById(ctx context.Context, id int32) (sqlc.GetVariantByIdRow, error) {
	return pr.DB.GetVariantById(ctx, id)
}

func (pr *productRepository) GetListProducts(ctx context.Context, arg sqlc.GetListProductsParams) ([]sqlc.GetListProductsRow, error) {
	return pr.DB.GetListProducts(ctx, arg)
}

func (pr *productRepository) CountGetListProducts(ctx context.Context, arg sqlc.CountGetListProductsParams) (int64, error) {
	return pr.DB.CountGetListProducts(ctx, arg)
}
