package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"nineshop-be/internal/db/sqlc"
)

type UserRepository interface {
	CreateUser(context context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	FindByEmail(context context.Context, email string) (sqlc.User, error)
	GetAllUser(context context.Context) ([]sqlc.User, error)
	GetAllUserV2(ctx context.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.GetAllUserV2Row, error)
	CountUser(ctx context.Context, search, role string, isActive *bool, page, limit int32) (int64, error)
	SoftDeleteUser(c context.Context, uuid pgtype.UUID) (sqlc.User, error)
	UpdateUser(context context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
	ActiveUser(ctx context.Context, userUuid pgtype.UUID) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, arg sqlc.CreateProductParams) (sqlc.Product, error)
	GetAllProductByFilter(ctx context.Context, limit, page, categoryId, minPrice, maxPrice, brandId int32, search, status string) ([]sqlc.GetAllProductByFilterRow, error)
	CountProduct(ctx context.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) (int64, error)
	GetProductById(ctx context.Context, id int32) (sqlc.Product, error)
	UpdateThumbnail(ctx context.Context, thumbnail string, id int32) (sqlc.Product, error)
	GetProductByCategoryId(ctx context.Context, id int32) ([]sqlc.GetProductByCategoryIdRow, error)
	GetProductBySlug(ctx context.Context, slug string) (sqlc.GetProductBySlugRow, error)
	GetTop8ProductSeller(ctx context.Context, cateID *int32) ([]sqlc.GetTop8ProductSellerRow, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, name string) (sqlc.Category, error)
	GetAll(ctx context.Context) ([]sqlc.Category, error)
}

type BrandRepository interface {
	GetAll(ctx context.Context) ([]sqlc.Brand, error)
}

type ProductImagesRepository interface {
	Save(ctx context.Context, arg sqlc.SaveAndUploadImgParams) (sqlc.ProductImage, error)
	GetImagesByProductId(ctx context.Context, id int32) ([]string, error)
}

type CartRepository interface {
	AddToCart(ctx context.Context, arg sqlc.AddToCartParams) (sqlc.Cart, error)
	ExistsProductId(ctx context.Context, arg sqlc.ExistsProductIdParams) (bool, error)
	UpdateCart(ctx context.Context, arg sqlc.UpdateCartParams) (sqlc.Cart, error)
	GetCartsByUserId(ctx context.Context, userID int32) ([]sqlc.GetCartsByUserIdRow, error)
	DeleteItemInCart(ctx context.Context, arg sqlc.DeleteItemInCartParams) error
	UpdateAllCart(ctx context.Context, arg sqlc.UpdateAllCartParams) (sqlc.Cart, error)
}

type PendingOrderItemRepository interface {
	Create(ctx context.Context, arg sqlc.CreatePendingOrderItemParams) (sqlc.PendingOrderItem, error)
	GetByPOrderItemId(ctx context.Context, id int32) ([]sqlc.PendingOrderItem, error)
	GetCountItem(ctx context.Context, id int32) (int64, error)
}

type PendingOrderRepository interface {
	Create(ctx context.Context, arg sqlc.CreatePendingOrderParams) (sqlc.PendingOrder, error)
	GetPOrderById(ctx context.Context, id int32) (sqlc.PendingOrder, error)
}

type PaymentMethodRepository interface {
	GetAllPayment(ctx context.Context) ([]sqlc.PaymentMethod, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, arg sqlc.CreateOrderParams) (sqlc.Order, error)
	UpdateStatusForUser(ctx context.Context, arg sqlc.UpdateStatusForUserParams) error
	GetOrderById(ctx context.Context, orderID int32) (sqlc.Order, error)
	GetAllOrders(ctx context.Context) ([]sqlc.GetAllOrdersRow, error)
	GetOrderDetailById(ctx context.Context, id int32) ([]sqlc.GetOrderDetailByIdRow, error)
	UpdateOrderPayment(ctx context.Context, arg sqlc.UpdateOrderPaymentParams) error
	GetListOrderByOrderId(ctx context.Context, arg sqlc.GetListOrderByOrderIdParams) ([]sqlc.GetListOrderByOrderIdRow, error)
	ViewDetailForMyOrder(ctx context.Context, orderID int32) ([]sqlc.ViewDetailForMyOrderRow, error)
}

type OrderItemRepository interface {
	AddOrderItem(ctx context.Context, arg sqlc.AddOrderItemParams) (sqlc.OrderItem, error)
	GetListOrderItemsByUserId(ctx context.Context, id int32) ([]sqlc.GetOrderItemByUserIdRow, error)
	GetCountItem(ctx context.Context, id int32) (int64, error)
}
type OrderStatusHistoryRepository interface {
	CreateOrderStatusHistory(ctx context.Context, arg sqlc.CreateOrderStatusHistoryParams) (sqlc.OrderStatusHistory, error)
	GetAllStatusByOrderId(ctx context.Context, orderID int32) ([]sqlc.GetAllStatusByOrderIdRow, error)
	GetAllStatusByOrderIdV2(ctx context.Context, id int32) ([]sqlc.GetAllStatusByOrderIdV2Row, error)
}
