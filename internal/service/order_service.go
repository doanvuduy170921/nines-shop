package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/constant"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"strings"
)

type orderService struct {
	repo repository.OrderRepository
	osh  repository.OrderStatusHistoryRepository
	ui   UserInject
}

func NewOrderService(repo repository.OrderRepository, osh repository.OrderStatusHistoryRepository, ui UserInject) OrderService {
	return &orderService{
		repo: repo,
		osh:  osh,
		ui:   ui,
	}
}

func (ns *orderService) UpdateStatusForUser(c *gin.Context, input dto.UpdateStatusForUserParams) error {
	ctx := c.Request.Context()
	// kiểm tra xem order có tồn tại hay kh
	order, err := ns.repo.GetOrderById(ctx, int32(input.OrderId))
	if err != nil {
		return utils.WrapError(err, "Order not found", http.StatusBadRequest)
	}
	// kiểm tra xem status truyền vào có hợp lệ kh
	if !constant.IsValidOrderStatus(strings.TrimSpace(input.Status)) {
		return utils.NewError(400, "invalid order status")
	}
	// kiểm tra tránh trường hợp nhảy cóc
	if !constant.CanTransaction(order.Status, input.Status) {
		return utils.NewError(400, "can not transaction for current status")
	}
	// update status order
	err = ns.repo.UpdateStatusForUser(ctx, sqlc.UpdateStatusForUserParams{
		Status:  input.Status,
		OrderID: int32(input.OrderId),
	})
	if err != nil {
		return utils.WrapError(err, "Update status for user fail", http.StatusBadRequest)
	}
	// cập nhật tracking order
	note := constant.OrderStatusNotes[input.Status]
	_, err = ns.osh.CreateOrderStatusHistory(ctx, sqlc.CreateOrderStatusHistoryParams{
		OrderID: int32(input.OrderId),
		Status:  input.Status,
		Note:    &note,
	})
	if err != nil {
		return utils.WrapError(err, "Create order status history fail", http.StatusBadRequest)
	}
	return ns.repo.UpdateStatusForUser(ctx, sqlc.UpdateStatusForUserParams{})
}

func (ns *orderService) GetAllOrders(c *gin.Context) ([]sqlc.GetAllOrdersRow, error) {
	ctx := c.Request.Context()
	orders, err := ns.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, utils.WrapError(err, "Get all orders fail", http.StatusBadRequest)
	}
	return orders, nil
}

func (ns *orderService) GetOrderDetailById(c *gin.Context, id int32) ([]sqlc.GetOrderDetailByIdRow, error) {
	ctx := c.Request.Context()

	_, err := ns.repo.GetOrderById(ctx, id)
	if err != nil {
		return nil, utils.WrapError(err, "order not found", http.StatusBadRequest)
	}
	order, err := ns.repo.GetOrderDetailById(ctx, id)
	if err != nil {
		return nil, utils.WrapError(err, "Get order detail fail", http.StatusBadRequest)
	}
	return order, nil
}

func (ns *orderService) GetAllStatusByOrderId(c *gin.Context, orderID int32) ([]sqlc.GetAllStatusByOrderIdRow, error) {
	ctx := c.Request.Context()
	status, err := ns.osh.GetAllStatusByOrderId(ctx, orderID)
	if err != nil {
		return nil, utils.WrapError(err, "Get all status by orderId fail", http.StatusBadRequest)
	}
	return status, nil
}

func (ns *orderService) GetAllStatusByOrderIdV2(c *gin.Context, orderID int32) ([]sqlc.GetAllStatusByOrderIdV2Row, error) {
	ctx := c.Request.Context()
	status, err := ns.osh.GetAllStatusByOrderIdV2(ctx, orderID)
	if err != nil {
		return nil, utils.WrapError(err, "Get all status by orderId fail", http.StatusBadRequest)
	}
	return status, nil
}

func (ns *orderService) GetListOrdersDetailByUserId(c *gin.Context, search *int32, status *string, limit int32, page int32) ([]sqlc.GetListOrderByOrderIdRow, error) {

	userUuid, exists := c.Get("user_uuid")
	if !exists {
		return []sqlc.GetListOrderByOrderIdRow{}, utils.NewError(400, "user uuid not found in request context")
	}
	userUuidStr := userUuid.(string)
	Uuid, err := utils.StringToPgUuid(userUuidStr)
	if err != nil {
		return []sqlc.GetListOrderByOrderIdRow{}, utils.NewError(400, "convert user uuid to uuid fail")
	}
	user, err := ns.ui.GetByUuid(c, Uuid)
	if err != nil {
		return []sqlc.GetListOrderByOrderIdRow{}, utils.NewError(400, "get user by uuid fail")
	}
	ctx := c.Request.Context()
	offset := (page - 1) * limit
	listOrders, err := ns.repo.GetListOrderByOrderId(ctx, sqlc.GetListOrderByOrderIdParams{
		Limit:  limit,
		Offset: offset,
		UserID: user.ID,
		Status: status,
		Search: search,
	})
	if err != nil {
		return []sqlc.GetListOrderByOrderIdRow{}, utils.WrapError(err, fmt.Sprintf("Get list order detail by userId fail :%v", err), http.StatusBadRequest)
	}
	return listOrders, nil
}

func (ns *orderService) ViewDetailForMyOrder(c *gin.Context, orderID int32) ([]sqlc.ViewDetailForMyOrderRow, error) {
	ctx := c.Request.Context()
	listOrders, err := ns.repo.ViewDetailForMyOrder(ctx, orderID)
	if err != nil {
		return nil, utils.WrapError(err, "View order detail fail", http.StatusBadRequest)
	}
	return listOrders, nil
}
