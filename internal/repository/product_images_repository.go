package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type productImagesRepository struct {
	DB *sqlc.Queries
}

func NewProductImagesRepository(DB *sqlc.Queries) ProductImagesRepository {
	return &productImagesRepository{DB: DB}
}

func (pr *productImagesRepository) Save(ctx context.Context, arg sqlc.SaveAndUploadImgParams) (sqlc.ProductImage, error) {
	images, err := pr.DB.SaveAndUploadImg(ctx, arg)
	if err != nil {
		return sqlc.ProductImage{}, err
	}
	return images, nil
}
