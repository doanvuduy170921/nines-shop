package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/cache"
	"time"
)

type userService struct {
	repo  repository.UserRepository
	redis cache.RedisCacheService
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		repo:  repo,
		redis: cache.NewRedisCacheService(redisClient),
	}
}

func (us *userService) CreateUser(c *gin.Context, input *sqlc.CreateUserParams) (sqlc.User, error) {
	context := c.Request.Context()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.WrapError(err, "Fail to generate from password", utils.ErrorCodeUnauthorized))
	}
	input.Password = string(hashPassword)
	user, err := us.repo.CreateUser(context, *input)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Fail to create user", utils.ErrorCodeInternalError)
	}
	if err := us.redis.Clear("get_all_users"); err == nil {
		log.Println("✅ redis clear all users")
	}
	return user, nil
}

func (us *userService) GetAllUser(c *gin.Context) ([]sqlc.User, error) {
	var cacheData struct {
		Users []sqlc.User `json:"users"`
	}
	if err := us.redis.Get("get_all_user", &cacheData); err == nil {
		log.Println("✅  Get all user from redis successfully")
		return cacheData.Users, nil
	}

	context := c.Request.Context()

	users, err := us.repo.GetAllUser(context)
	if err != nil {
		return nil, utils.WrapError(err, "Fail to get all user", utils.ErrorCodeInternalError)
	}
	cacheData.Users = users
	if err := us.redis.Set("get_all_user", cacheData, 5*time.Minute); err != nil {
		log.Printf("Fail to set user to cache : %s", err.Error())
	}

	return users, nil
}

func (us *userService) GetAllUserV2(ctx *gin.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.User, int64, error) {
	context := ctx.Request.Context()

	users, err := us.repo.GetAllUserV2(context, search, role, isActive, page, limit)
	if err != nil {
		return nil, 0, utils.WrapError(err, "Fail to get all user V2", utils.ErrorCodeInternalError)
	}
	count, err := us.repo.CountUser(context, search, role, isActive, page, limit)
	if err != nil {
		return nil, 0, utils.WrapError(err, "Fail to count user V2", utils.ErrorCodeInternalError)
	}
	return users, count, nil
}

func (us *userService) SoftDeleteUser(c *gin.Context, uuid pgtype.UUID) (sqlc.User, error) {
	ctx := c.Request.Context()
	user, err := us.repo.SoftDeleteUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Fail to soft delete user", utils.ErrorCodeInternalError)
	}
	return user, nil
}

func (us *userService) UpdateUser(c *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := c.Request.Context()

	user, err := us.repo.UpdateUser(context, input)
	if err != nil {
		log.Printf("err : %v", err.Error())
		return sqlc.User{}, utils.WrapError(err, "Fail to update user", utils.ErrorCodeInternalError)
	}
	return user, nil
}

func (us *userService) GetByUuid(ctx *gin.Context, uuid pgtype.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	user, err := us.repo.GetByUuid(context, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return sqlc.User{}, utils.WrapError(err, "Not find user with uuid:"+uuid.String(), utils.ErrorCodeUnauthorized)
		}
		return sqlc.User{}, utils.WrapError(err, "Fail to get user by uuid", utils.ErrorCodeInternalError)
	}
	return user, nil
}
