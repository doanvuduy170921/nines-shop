package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type brandRepository struct {
	DB *sqlc.Queries
}

func NewBrandRepository(DB *sqlc.Queries) BrandRepository {
	return &brandRepository{DB: DB}
}

func (bu *brandRepository) GetAll(ctx context.Context) ([]sqlc.Brand, error) {
	brands, err := bu.DB.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return brands, nil
}
