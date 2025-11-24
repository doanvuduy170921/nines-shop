package repository

import (
	"context"

	"nineshop-be/internal/db/sqlc"
)

type paymentMethodRepository struct {
	DB *sqlc.Queries
}

func NewPaymentMethodRepository(db *sqlc.Queries) PaymentMethodRepository {
	return &paymentMethodRepository{DB: db}
}

func (pmr *paymentMethodRepository) GetAllPayment(ctx context.Context) ([]sqlc.PaymentMethod, error) {
	return pmr.DB.GetAllPayment(ctx)
}
