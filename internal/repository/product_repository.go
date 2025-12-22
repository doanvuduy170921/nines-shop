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

func (pr *productRepository) CreateProduct(ctx context.Context, arg sqlc.CreateProductParams) (sqlc.Product, error) {
	product, err := pr.DB.CreateProduct(ctx, arg)
	if err != nil {
		return sqlc.Product{}, err
	}
	return product, nil
}

func (pr *productRepository) GetAllProductByFilter(ctx context.Context, limit, page, categoryId, minPrice, maxPrice, brandId int32, search, status string) ([]sqlc.GetAllProductByFilterRow, error) {
	offset := (page - 1) * limit

	products, err := pr.DB.GetAllProductByFilter(ctx, sqlc.GetAllProductByFilterParams{
		Limit:      limit,
		Offset:     offset,
		Search:     search,
		Status:     status,
		CategoryID: categoryId,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		BrandID:    brandId,
	})
	if err != nil {
		return nil, err
	}
	return products, nil
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

func (pr *productRepository) GetProductByCategoryId(ctx context.Context, id int32) ([]sqlc.GetProductByCategoryIdRow, error) {
	products, err := pr.DB.GetProductByCategoryId(ctx, id)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (pr *productRepository) GetProductBySlug(ctx context.Context, slug string) (sqlc.GetProductBySlugRow, error) {
	product, err := pr.DB.GetProductBySlug(ctx, slug)
	if err != nil {
		return sqlc.GetProductBySlugRow{}, err
	}
	return product, nil
}

func (pr *productRepository) GetTop8ProductSeller(ctx context.Context, cateID *int32) ([]sqlc.GetTop8ProductSellerRow, error) {
	return pr.DB.GetTop8ProductSeller(ctx, cateID)
}
