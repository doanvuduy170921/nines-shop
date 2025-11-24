package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type orderStatusHistoryRepository struct {
	DB *sqlc.Queries
}

func NewOrderStatusHistory(db *sqlc.Queries) OrderStatusHistoryRepository {
	return &orderStatusHistoryRepository{DB: db}
}

func (nh *orderStatusHistoryRepository) CreateOrderStatusHistory(ctx context.Context, arg sqlc.CreateOrderStatusHistoryParams) (sqlc.OrderStatusHistory, error) {
	return nh.DB.CreateOrderStatusHistory(ctx, arg)
}
