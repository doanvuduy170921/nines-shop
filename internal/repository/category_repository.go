package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type categoryRepository struct {
	DB *sqlc.Queries
}

func NewCategoryRepository(DB *sqlc.Queries) CategoryRepository {
	return &categoryRepository{
		DB: DB,
	}
}

func (cr *categoryRepository) Create(ctx context.Context, name string) (sqlc.Category, error) {
	category, err := cr.DB.CreateCategory(ctx, name)
	if err != nil {
		return sqlc.Category{}, err
	}
	return category, nil
}

func (cr *categoryRepository) GetAll(ctx context.Context) ([]sqlc.Category, error) {
	categories, err := cr.DB.GetAllCategories(ctx)
	if err != nil {
		return []sqlc.Category{}, err
	}
	return categories, nil
}
