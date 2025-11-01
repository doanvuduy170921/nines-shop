package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
)

type cartService struct {
	repo repository.CartRepository
	uu   UserUpdater
}
type UserUpdater interface {
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
}

func NewCartService(repo repository.CartRepository, uu UserUpdater) CartService {
	return &cartService{
		repo: repo,
		uu:   uu}
}

func (cs *cartService) AddToCart(ctx *gin.Context, userUuid string, req dto.AddToCartParams) (sqlc.Cart, error) {
	c := ctx.Request.Context()
	body := dto.MapParamToSqlcCart(req)
	uuid, err := utils.StringToPgUuid(userUuid)
	if err != nil {
		return sqlc.Cart{}, utils.WrapError(err, "Convert fail string to pgUUid", http.StatusBadRequest)
	}
	user, err := cs.uu.GetByUuid(c, uuid)
	if err != nil {
		return sqlc.Cart{}, utils.WrapError(err, "Get user fail by uuid", http.StatusBadRequest)
	}
	// kiểm tra xem đã có sp trong giỏ hay chưa, nếu có update quantity , chưa thì add to cart
	exists, err := cs.repo.ExistsProductId(ctx, sqlc.ExistsProductIdParams{
		ProductID: body.ProductID,
		UserID:    user.ID,
	})
	if err != nil {
		return sqlc.Cart{}, utils.WrapError(err, "Check fail exists productId in cart ", http.StatusBadRequest)
	}
	var cart sqlc.Cart
	if !exists {
		cart, err = cs.repo.AddToCart(c, sqlc.AddToCartParams{
			UserID:    int32(user.ID),
			ProductID: body.ProductID,
			Quantity:  body.Quantity,
		})
		if err != nil {
			return sqlc.Cart{}, utils.WrapError(err, "Add to cart fail", http.StatusBadRequest)
		}
	} else {
		cart, err = cs.repo.UpdateCart(ctx, sqlc.UpdateCartParams{
			ProductID: body.ProductID,
			UserID:    int32(user.ID),
			Quantity:  *body.Quantity,
		})
		if err != nil {
			return sqlc.Cart{}, utils.WrapError(err, "Update cart fail", http.StatusBadRequest)
		}
	}

	return cart, nil
}
