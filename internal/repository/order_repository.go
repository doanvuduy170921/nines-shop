package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type orderRepository struct {
	DB *sqlc.Queries
}

func NewOrderRepository(db *sqlc.Queries) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (nr *orderRepository) CreateOrder(ctx context.Context, arg sqlc.CreateOrderParams) (sqlc.Order, error) {
	order, err := nr.DB.CreateOrder(ctx, arg)
	if err != nil {
		return sqlc.Order{}, err
	}
	return order, nil
}
