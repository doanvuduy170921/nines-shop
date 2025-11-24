package service

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/db/sqlc"

	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
)

type orderItemService struct {
	repo repository.OrderItemRepository
	ui   UserInject
}

func NewOrderItemService(repo repository.OrderItemRepository, ui UserInject) OrderItemService {
	return &orderItemService{
		repo: repo,
		ui:   ui,
	}
}

func (nh *orderItemService) GetListOrderItemByUserId(ctx *gin.Context) ([]sqlc.GetOrderItemByUserIdRow, error) {
	c := ctx.Request.Context()
	userUuid, exists := ctx.Get("user_uuid")
	if !exists {
		return []sqlc.GetOrderItemByUserIdRow{}, utils.NewError(400, "user uuid not found in request context")
	}
	userUuidStr := userUuid.(string)
	Uuid, err := utils.StringToPgUuid(userUuidStr)
	if err != nil {
		return []sqlc.GetOrderItemByUserIdRow{}, utils.NewError(400, "convert user uuid to uuid fail")
	}
	user, err := nh.ui.GetByUuid(ctx, Uuid)
	if err != nil {
		return []sqlc.GetOrderItemByUserIdRow{}, utils.NewError(400, "get user by uuid fail")
	}

	orderItems, err := nh.repo.GetListOrderItemsByUserId(c, user.ID)
	if err != nil {
		return []sqlc.GetOrderItemByUserIdRow{}, utils.NewError(400, "get get List order items by userId fail")
	}
	return orderItems, nil
}
