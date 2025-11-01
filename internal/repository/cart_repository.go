package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type cartRepository struct {
	DB *sqlc.Queries
}

func NewCartRepository(DB *sqlc.Queries) CartRepository {
	return &cartRepository{DB: DB}
}

func (cs *cartRepository) AddToCart(ctx context.Context, arg sqlc.AddToCartParams) (sqlc.Cart, error) {
	cart, err := cs.DB.AddToCart(ctx, arg)
	return cart, err
}

func (cs *cartRepository) ExistsProductId(ctx context.Context, arg sqlc.ExistsProductIdParams) (bool, error) {
	exist, err := cs.DB.ExistsProductId(ctx, arg)
	return exist, err
}

func (cs *cartRepository) UpdateCart(ctx context.Context, arg sqlc.UpdateCartParams) (sqlc.Cart, error) {
	cart, err := cs.DB.UpdateCart(ctx, arg)
	if err != nil {
		return sqlc.Cart{}, err
	}
	return cart, nil
}
