package repository

import (
	"context"
	"fmt"
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

func (pr *productRepository) GetAllProductByFilter(ctx context.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) ([]sqlc.Product, error) {
	offset := (page - 1) * limit

	fmt.Printf("Params: limit=%d, offset=%d, search=%s, status=%s, categoryId=%d, minPrice=%d, maxPrice=%d\n",
		limit, offset, search, status, categoryId, minPrice, maxPrice)
	products, err := pr.DB.GetAllProductByFilter(ctx, sqlc.GetAllProductByFilterParams{
		Limit:      limit,
		Offset:     offset,
		Search:     search,
		Status:     status,
		CategoryID: categoryId,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
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
