package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
)

type UserService interface {
	CreateUser(c *gin.Context, input *sqlc.CreateUserParams) (sqlc.User, error)
	GetAllUser(c *gin.Context) ([]sqlc.User, error)
	GetAllUserV2(ctx *gin.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.User, int64, error)
	SoftDeleteUser(c *gin.Context, uuid pgtype.UUID) (sqlc.User, error)
	UpdateUser(c *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	GetByUuid(ctx *gin.Context, uuid pgtype.UUID) (sqlc.User, error)
}

type AuthService interface {
	Login(c *gin.Context, email, password string) (string, string, int, string, error)
	RefreshToken(c *gin.Context, tokenString string) (string, string, int, error)
	Logout(c *gin.Context, refreshTokenString string) error
}

type ProductService interface {
	CreateProduct(ctx *gin.Context, arg dto.CreateProductParamDto) (sqlc.Product, error)
	GetAllByFilter(ctx *gin.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) ([]sqlc.GetAllProductByFilterRow, int64, error)
	GetProductByCategoryId(ctx *gin.Context, id int32) ([]sqlc.GetProductByCategoryIdRow, error)
	GetProductBySlug(ctx *gin.Context, slug string) (sqlc.GetProductBySlugRow, error)
	GetImagesBySlug(c *gin.Context, slug string) ([]string, error)
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
	Create(ctx *gin.Context, arg dto.CreatePendingOrderDto) (sqlc.PendingOrder, error)
	ValidateOTP(ctx *gin.Context, arg dto.ValidateOTPParams) (sqlc.Order, error)
}

type PaymentService interface {
	GetAllPayment(ctx *gin.Context) ([]sqlc.PaymentMethod, error)
}

type OrderItemService interface {
	GetListOrderItemByUserId(ctx *gin.Context) ([]sqlc.GetOrderItemByUserIdRow, error)
}
