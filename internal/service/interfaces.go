package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/url"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
)

type UserService interface {
	CreateUser(c *gin.Context, input *sqlc.CreateUserParams) (dto.CreateUserRes, error)
	GetAllUser(c *gin.Context) ([]sqlc.User, error)
	GetAllUserV2(ctx *gin.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.GetAllUserV2Row, int64, error)
	SoftDeleteUser(c *gin.Context, uuid pgtype.UUID) (sqlc.User, error)
	UpdateUser(c *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	GetByUuid(ctx *gin.Context, uuid pgtype.UUID) (sqlc.User, error)
	ActiveUser(ctx context.Context, email string, otp string) error
}

type AuthService interface {
	Login(c *gin.Context, email, password string) (string, string, int, string, error)
	RefreshToken(c *gin.Context, tokenString string) (string, string, int, error)
	Logout(c *gin.Context, refreshTokenString string) error
}

type ProductService interface {
	CreateProduct(ctx *gin.Context, arg dto.CreateProductParamDto) (sqlc.Product, error)
	GetAllByFilter(ctx *gin.Context, limit, page, categoryId, minPrice, maxPrice, brandId int32, search, status string) ([]sqlc.GetAllProductByFilterRow, int64, error)
	GetProductByCategoryId(ctx *gin.Context, id int32) ([]sqlc.GetProductByCategoryIdRow, error)
	GetProductBySlug(ctx *gin.Context, slug string) (sqlc.GetProductBySlugRow, error)
	GetImagesBySlug(c *gin.Context, slug string) ([]string, error)
	GetTop8ProductSeller(ctx *gin.Context, cateID *int32) ([]sqlc.GetTop8ProductSellerRow, error)
}

type CategoryService interface {
	Create(ctx *gin.Context, name string) (sqlc.Category, error)
	GetAll(ctx *gin.Context) ([]sqlc.Category, error)
}

type BrandService interface {
	GetAll(ctx *gin.Context) ([]sqlc.Brand, error)
}

type ProductImagesService interface {
	SaveUploadFile(ctx *gin.Context, id int64) (*dto.UploadResult, error)
}

type CartService interface {
	AddToCart(ctx *gin.Context, userUuid string, req dto.AddToCartParams) (sqlc.Cart, error)
	GetCartsByUserId(ctx *gin.Context, userUuid string) ([]sqlc.GetCartsByUserIdRow, error)
	DeleteItem(ctx *gin.Context, input dto.DeleteItemInCartParams) error
	UpdateAllCart(ctx *gin.Context, input dto.UpdateAllCartParam) ([]sqlc.Cart, float64, float64, float64, float64, error)
}

type PendingOrderService interface {
	Create(ctx *gin.Context, arg dto.CreatePendingOrderDto) (dto.CreatePendingOrderResponse, error)
	ValidateOTP(ctx *gin.Context, arg dto.ValidateOTPParams) (sqlc.Order, error)
}

type PaymentService interface {
	GetAllPayment(ctx *gin.Context) ([]sqlc.PaymentMethod, error)
	HandleVnPaySuccess(ctx *gin.Context, orderRef string, query url.Values) error
	ProcessVNPayPayment(ctx *gin.Context, orderID int64, amount int64, transactionNo string, bankCode string) (sqlc.Order, error)
}

type OrderItemService interface {
	GetListOrderItemByUserId(ctx *gin.Context) ([]sqlc.GetOrderItemByUserIdRow, error)
}

type OrderService interface {
	UpdateStatusForUser(c *gin.Context, input dto.UpdateStatusForUserParams) error
	GetAllOrders(ctx *gin.Context) ([]sqlc.GetAllOrdersRow, error)
	GetOrderDetailById(c *gin.Context, id int32) ([]sqlc.GetOrderDetailByIdRow, error)
	GetAllStatusByOrderId(c *gin.Context, orderID int32) ([]sqlc.GetAllStatusByOrderIdRow, error)
	GetAllStatusByOrderIdV2(c *gin.Context, orderID int32) ([]sqlc.GetAllStatusByOrderIdV2Row, error)
	GetListOrdersDetailByUserId(c *gin.Context, search *int32, status *string, limit int32, page int32) ([]sqlc.GetListOrderByOrderIdRow, error)
	ViewDetailForMyOrder(c *gin.Context, orderID int32) ([]sqlc.ViewDetailForMyOrderRow, error)
}
