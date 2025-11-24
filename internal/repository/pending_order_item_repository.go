package repository

import (
	"context"
	"nineshop-be/internal/db/sqlc"
)

type pendingOrderItemRepository struct {
	DB *sqlc.Queries
}

func NewPendingOrderItemRepository(db *sqlc.Queries) PendingOrderItemRepository {
	return &pendingOrderItemRepository{
		DB: db,
	}
}

func (nr *pendingOrderItemRepository) Create(ctx context.Context, arg sqlc.CreatePendingOrderItemParams) (sqlc.PendingOrderItem, error) {
	pendingOrder, err := nr.DB.CreatePendingOrderItem(ctx, arg)
	if err != nil {
		return sqlc.PendingOrderItem{}, err
	}
	return pendingOrder, nil
}

func (nr *pendingOrderItemRepository) GetByPOrderItemId(ctx context.Context, id int32) ([]sqlc.PendingOrderItem, error) {
	return nr.DB.GetByPOrderItemId(ctx, id)
}
