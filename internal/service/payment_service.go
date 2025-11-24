package service

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
)

type paymentService struct {
	repo repository.PaymentMethodRepository
}

func NewPaymentService(repo repository.PaymentMethodRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}

func (ps *paymentService) GetAllPayment(ctx *gin.Context) ([]sqlc.PaymentMethod, error) {
	c := ctx.Request.Context()
	return ps.repo.GetAllPayment(c)
}
