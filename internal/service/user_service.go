package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/utils"
	"nineshop-be/pkg/cache"
	"nineshop-be/pkg/email"
	"time"
)

type userService struct {
	repo      repository.UserRepository
	redis     cache.RedisCacheService
	mail      email.EmailService
	mailQueue chan mailPayload
}

type mailPayload struct {
	Receiver string
	OTP      string
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client, mail email.EmailService) UserService {
	// 1. Khởi tạo instance của userService trước
	us := &userService{
		repo:      repo,
		redis:     cache.NewRedisCacheService(redisClient),
		mail:      mail,
		mailQueue: make(chan mailPayload, 100), // Đảm bảo channel đã được tạo
	}

	// 2. Chạy worker từ instance 'us' vừa tạo
	go us.ProcessMailQueue()

	return us
}
func (us *userService) CreateUser(c *gin.Context, input *sqlc.CreateUserParams) (dto.CreateUserRes, error) {
	ctx := c.Request.Context()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.CreateUserRes{}, utils.WrapError(err, "Fail to generate from password", utils.ErrorCodeUnauthorized)
	}
	input.Password = string(hashPassword)
	user, err := us.repo.CreateUser(ctx, *input)
	if err != nil {
		return dto.CreateUserRes{}, utils.WrapError(err, "Fail to create user", utils.ErrorCodeInternalError)
	}
	otp := utils.GenOTP()
	key := cache.OTPCreate(input.Email)

	expiredAt := 5 * time.Minute
	duration := time.Now().Add(expiredAt)
	if err := us.redis.Set(key, otp, expiredAt); err != nil {
		return dto.CreateUserRes{}, utils.WrapError(err, "Fail to set otp", utils.ErrorCodeInternalError)
	}
	select {
	case us.mailQueue <- mailPayload{
		Receiver: user.Email,
		OTP:      otp,
	}:
	default:
		log.Println("Email queue is full", user.Email)
	}
	userRes := dto.CreateUserRes{
		Email:     user.Email,
		ExpiredAt: duration.Format("2006-01-02 15:04:05"),
	}
	return userRes, nil
}

func (us *userService) GetAllUser(c *gin.Context) ([]sqlc.User, error) {
	var cacheData struct {
		Users []sqlc.User `json:"users"`
	}
	if err := us.redis.Get("get_all_user", &cacheData); err == nil {
		log.Println("✅ Get all user from redis successfully")
		return cacheData.Users, nil
	}

	ctx := c.Request.Context()

	users, err := us.repo.GetAllUser(ctx)
	if err != nil {
		return nil, utils.WrapError(err, "Fail to get all user", utils.ErrorCodeInternalError)
	}
	cacheData.Users = users
	if err := us.redis.Set("get_all_user", cacheData, 5*time.Minute); err != nil {
		log.Printf("Fail to set user to cache : %s", err.Error())
	}

	return users, nil
}

func (us *userService) GetAllUserV2(c *gin.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.GetAllUserV2Row, int64, error) {
	ctx := c.Request.Context()

	users, err := us.repo.GetAllUserV2(ctx, search, role, isActive, page, limit)
	if err != nil {
		return nil, 0, utils.WrapError(err, "Fail to get all user V2", utils.ErrorCodeInternalError)
	}
	count, err := us.repo.CountUser(ctx, search, role, isActive, page, limit)
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
	ctx := c.Request.Context()

	user, err := us.repo.UpdateUser(ctx, input)
	if err != nil {
		log.Printf("err : %v", err.Error())
		return sqlc.User{}, utils.WrapError(err, "Fail to update user", utils.ErrorCodeInternalError)
	}
	return user, nil
}

func (us *userService) GetByUuid(c *gin.Context, uuid pgtype.UUID) (sqlc.User, error) {
	ctx := c.Request.Context()
	user, err := us.repo.GetByUuid(ctx, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return sqlc.User{}, utils.WrapError(err, "Not find user with uuid:"+uuid.String(), utils.ErrorCodeUnauthorized)
	}
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "Fail to get user by uuid", utils.ErrorCodeInternalError)
	}
	return user, nil
}

func (us *userService) ActiveUser(ctx context.Context, email string, otp string) error {

	user, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		return utils.WrapError(err, "Fail to find user by email", utils.ErrorCodeInternalError)
	}
	if user.IsActive != nil && *user.IsActive {
		return utils.NewError(http.StatusBadRequest, "User is already active")
	}
	var redisOTP string
	key := cache.OTPCreate(email)
	if err := us.redis.Get(key, &redisOTP); err != nil {
		if errors.Is(err, redis.Nil) {
			return utils.NewError(http.StatusBadRequest, "OTP expired or invalid")
		}
		return utils.WrapError(err, "Fail to get otp", http.StatusInternalServerError)
	}
	if otp != redisOTP {
		return utils.NewError(http.StatusBadRequest, "OTP is not valid")
	}
	if err := us.repo.ActiveUser(ctx, user.UserUuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError(http.StatusNotFound, "User not found")
		}
		return utils.WrapError(err, "Fail to active user", http.StatusInternalServerError)
	}
	if err := us.redis.Clear(key); err != nil {
		log.Printf("Fail to clear cache : %s", err.Error())
	}
	return nil
}

func (us *userService) ProcessMailQueue() {
	for payload := range us.mailQueue {
		err := us.mail.SendOTPEmail(payload.Receiver, payload.OTP)
		if err != nil {
			log.Printf("Error sending otp email to %s: %s", payload.Receiver, err.Error())
		} else {
			log.Printf("Successfully send otp email to %s", payload.Receiver)
		}
	}
}
