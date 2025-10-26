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
	GetAllByFilter(ctx *gin.Context, limit, page, categoryId, minPrice, maxPrice int32, search, status string) ([]sqlc.Product, int64, error)
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
