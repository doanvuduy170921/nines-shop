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
func (nh *orderStatusHistoryRepository) GetAllStatusByOrderId(ctx context.Context, orderID int32) ([]sqlc.GetAllStatusByOrderIdRow, error) {
	return nh.DB.GetAllStatusByOrderId(ctx, orderID)
}
func (nh *orderStatusHistoryRepository) GetAllStatusByOrderIdV2(ctx context.Context, id int32) ([]sqlc.GetAllStatusByOrderIdV2Row, error) {
	return nh.DB.GetAllStatusByOrderIdV2(ctx, id)
}
