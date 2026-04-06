package service

import (
	"context"
	"errors"
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
	pu   ProductUpdater
}
type UserUpdater interface {
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
}

func NewCartService(repo repository.CartRepository, uu UserUpdater, pu ProductUpdater) CartService {
	return &cartService{
		repo: repo,
		uu:   uu,
		pu:   pu}
}

var shippingOptions = map[string]float64{
	"Standard Delivery": 4.99,
	"Express Delivery":  12.99,
	"Free Shipping":     0.00,
}

func (cs *cartService) AddToCart(ctx *gin.Context, userUuid string, req dto.AddToCartParams) (sqlc.Cart, error) {
	c := ctx.Request.Context()
	body := dto.MapParamToSqlcCart(req)
	if req.Quantity <= 0 {
		return sqlc.Cart{}, utils.NewError(http.StatusBadRequest, "Cart's quantity must be positive")
	}

	variant, err := cs.pu.GetVariantById(ctx, body.VariantID)
	if err != nil {
		return sqlc.Cart{}, utils.WrapError(err, "Variant of product not found", http.StatusBadRequest)
	}
	if int32(req.Quantity) > variant.StockQuantity {
		return sqlc.Cart{}, utils.WrapError(err, "Quantity is too high", http.StatusBadRequest)
	}
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
		VariantID: body.VariantID,
		UserID:    user.ID,
	})
	if err != nil {
		return sqlc.Cart{}, utils.WrapError(err, "Check fail exists productId in cart ", http.StatusBadRequest)
	}
	var cart sqlc.Cart
	if !exists {
		cart, err = cs.repo.AddToCart(c, sqlc.AddToCartParams{
			UserID:    int32(user.ID),
			VariantID: body.VariantID,
			Quantity:  body.Quantity,
		})
		if err != nil {
			return sqlc.Cart{}, utils.WrapError(err, "Add to cart fail", http.StatusBadRequest)
		}
	} else {
		cart, err = cs.repo.UpdateCart(ctx, sqlc.UpdateCartParams{
			VariantID: body.VariantID,
			UserID:    user.ID,
			Quantity:  *body.Quantity,
		})
		if err != nil {
			return sqlc.Cart{}, utils.WrapError(err, "Update cart fail", http.StatusBadRequest)
		}
	}

	return cart, nil
}

func (cs *cartService) GetCartsByUserId(ctx *gin.Context, userUuid string) ([]sqlc.GetCartsByUserIdRow, error) {
	c := ctx.Request.Context()

	uuid, err := utils.StringToPgUuid(userUuid)
	if err != nil {
		return nil, utils.WrapError(err, "Convert fail string to pgUUid", http.StatusBadRequest)
	}
	user, err := cs.uu.GetByUuid(c, uuid)
	if err != nil {
		return nil, utils.WrapError(err, "Get user fail by uuid", http.StatusBadRequest)
	}

	carts, err := cs.repo.GetCartsByUserId(c, user.ID)
	if err != nil {
		return []sqlc.GetCartsByUserIdRow{}, err
	}
	return carts, nil
}

func (cs *cartService) DeleteItem(ctx *gin.Context, input dto.DeleteItemInCartParams) error {
	c := ctx.Request.Context()
	userUuidStr, exists := ctx.Get("user_uuid")
	if !exists {
		return errors.New("User uuid not found in context")
	}
	userUuid := userUuidStr.(string)
	uuid, err := utils.StringToPgUuid(userUuid)
	user, err := cs.uu.GetByUuid(c, uuid)
	if err != nil {
		return utils.WrapError(err, "Get user fail by uuid", http.StatusBadRequest)
	}

	err = cs.repo.DeleteItemInCart(c, sqlc.DeleteItemInCartParams{
		UserID:    user.ID,
		VariantID: input.VariantID,
	})
	if err != nil {
		return utils.WrapError(err, "Delete item fail", http.StatusBadRequest)
	}
	return nil
}

func (cs *cartService) UpdateAllCart(ctx *gin.Context, input dto.UpdateAllCartParam) ([]sqlc.Cart, float64, float64, float64, float64, error) {
	c := ctx.Request.Context()
	userUuidStr, exists := ctx.Get("user_uuid")
	if !exists {
		return []sqlc.Cart{}, 0, 0, 0, 0, errors.New("User uuid not found in context")
	}
	userUuid := userUuidStr.(string)
	uuid, err := utils.StringToPgUuid(userUuid)
	user, err := cs.uu.GetByUuid(c, uuid)
	if err != nil {
		return []sqlc.Cart{}, 0, 0, 0, 0, utils.WrapError(err, "Get user fail by uuid", http.StatusBadRequest)
	}
	var cart []sqlc.Cart

	var subtotal float64
	for _, item := range input.Items {
		if exists, err = cs.repo.ExistsProductId(ctx, sqlc.ExistsProductIdParams{
			VariantID: int32(item.VariantID),
			UserID:    user.ID,
		}); exists && err == nil {
			cartItem, err := cs.repo.UpdateAllCart(ctx, sqlc.UpdateAllCartParams{
				VariantID: int32(item.VariantID),
				UserID:    user.ID,
				Quantity:  int32(item.Quantity),
			})
			if err != nil {
				return []sqlc.Cart{}, 0, 0, 0, 0, utils.WrapError(err, "Update cart fail", http.StatusBadRequest)
			}
			cart = append(cart, cartItem)
			subtotal += float64(item.Quantity) * item.Price
		}
	}
	var tax float64 = subtotal * 0.01
	shipPrice, ok := shippingOptions[input.Shipping]
	if !ok {
		return []sqlc.Cart{}, 0, 0, 0, 0, utils.WrapError(err, "Shipping method is invalid", http.StatusBadRequest)
	}
	total := subtotal + tax + shipPrice
	return cart, subtotal, total, tax, shipPrice, nil
}
