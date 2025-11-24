package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type pendingOrderRepository struct {
	DB *sqlc.Queries
}

func NewPendingOrderRepository(db *sqlc.Queries) PendingOrderRepository {
	return &pendingOrderRepository{
		DB: db,
	}
}

func (nr *pendingOrderRepository) Create(ctx context.Context, arg sqlc.CreatePendingOrderParams) (sqlc.PendingOrder, error) {
	pendingOrder, err := nr.DB.CreatePendingOrder(ctx, arg)
	if err != nil {
		return sqlc.PendingOrder{}, err
	}
	return pendingOrder, nil
}

func (nr *pendingOrderRepository) GetPOrderById(ctx context.Context, id int32) (sqlc.PendingOrder, error) {
	return nr.DB.GetPOrderById(ctx, id)
}
