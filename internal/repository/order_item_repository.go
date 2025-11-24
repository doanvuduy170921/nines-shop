package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type orderItemRepository struct {
	DB *sqlc.Queries
}

func NewOrderItemRepository(db *sqlc.Queries) OrderItemRepository {
	return &orderItemRepository{
		DB: db,
	}
}

func (or *orderItemRepository) AddOrderItem(ctx context.Context, arg sqlc.AddOrderItemParams) (sqlc.OrderItem, error) {
	return or.DB.AddOrderItem(ctx, arg)
}

func (or *orderItemRepository) GetListOrderItemsByUserId(ctx context.Context, id int32) ([]sqlc.GetOrderItemByUserIdRow, error) {
	return or.DB.GetOrderItemByUserId(ctx, id)
}
