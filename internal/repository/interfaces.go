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
	GetAllUserV2(ctx context.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.User, error)
	CountUser(ctx context.Context, search, role string, isActive *bool, page, limit int32) (int64, error)
	SoftDeleteUser(c context.Context, uuid pgtype.UUID) (sqlc.User, error)
	UpdateUser(context context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error)
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, arg sqlc.CreateProductParams) (sqlc.Product, error)
	GetAllProductByFilter(ctx context.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) ([]sqlc.Product, error)
	CountProduct(ctx context.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) (int64, error)
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
}
