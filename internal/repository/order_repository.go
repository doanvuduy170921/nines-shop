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

func (nr *orderRepository) UpdateStatusForUser(ctx context.Context, arg sqlc.UpdateStatusForUserParams) error {
	return nr.DB.UpdateStatusForUser(ctx, arg)
}
func (nr *orderRepository) GetOrderById(ctx context.Context, orderID int32) (sqlc.Order, error) {
	return nr.DB.GetOrderById(ctx, orderID)
}

func (nr *orderRepository) GetAllOrders(ctx context.Context) ([]sqlc.GetAllOrdersRow, error) {
	return nr.DB.GetAllOrders(ctx)
}

func (nr *orderRepository) GetOrderDetailById(ctx context.Context, id int32) ([]sqlc.GetOrderDetailByIdRow, error) {
	return nr.DB.GetOrderDetailById(ctx, id)
}

func (nr *orderRepository) UpdateOrderPayment(ctx context.Context, arg sqlc.UpdateOrderPaymentParams) error {
	return nr.DB.UpdateOrderPayment(ctx, arg)
}

func (nr *orderRepository) GetListOrderByOrderId(ctx context.Context, arg sqlc.GetListOrderByOrderIdParams) ([]sqlc.GetListOrderByOrderIdRow, error) {
	return nr.DB.GetListOrderByOrderId(ctx, arg)
}

func (nr *orderRepository) ViewDetailForMyOrder(ctx context.Context, orderID int32) ([]sqlc.ViewDetailForMyOrderRow, error) {
	return nr.DB.ViewDetailForMyOrder(ctx, orderID)
}
