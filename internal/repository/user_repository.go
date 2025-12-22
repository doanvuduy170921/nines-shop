package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"nineshop-be/internal/db/sqlc"
)

type userRepository struct {
	DB *sqlc.Queries
}

func NewUserRepository(DB *sqlc.Queries) UserRepository {
	return &userRepository{
		DB: DB,
	}
}

func (ur *userRepository) CreateUser(context context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.DB.CreateUser(context, arg)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}
func (ur *userRepository) FindByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return ur.DB.FindByEmail(ctx, email)
}

func (ur *userRepository) GetAllUser(context context.Context) ([]sqlc.User, error) {
	return ur.DB.GetAllUser(context)
}

func (ur *userRepository) GetAllUserV2(ctx context.Context, search, role string, isActive *bool, page, limit int32) ([]sqlc.GetAllUserV2Row, error) {
	isActiveFilter := ""
	var isActiveBool bool

	if isActive != nil {
		isActiveFilter = "set"
		isActiveBool = *isActive
	}
	offset := (page - 1) * limit

	users, err := ur.DB.GetAllUserV2(ctx, sqlc.GetAllUserV2Params{
		Limit:          limit,
		Offset:         offset,
		Search:         &search,
		Role:           &role,
		IsActiveFilter: &isActiveFilter,
		IsActive:       &isActiveBool,
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *userRepository) CountUser(ctx context.Context, search, role string, isActive *bool, page, limit int32) (int64, error) {
	isActiveFilter := ""
	var isActiveBool bool

	if isActive != nil {
		isActiveFilter = "set"
		isActiveBool = *isActive
	}

	count, err := ur.DB.CountUser(ctx, sqlc.CountUserParams{
		Search:         &search,
		Role:           &role,
		IsActiveFilter: &isActiveFilter,
		IsActive:       &isActiveBool,
	})
	if err != nil {
		log.Printf("error : %v", err)
		return 0, err
	}
	return count, nil
}

func (ur *userRepository) SoftDeleteUser(c context.Context, uuid pgtype.UUID) (sqlc.User, error) {
	user, err := ur.DB.SoftDeleteUser(c, uuid)
	if err != nil {
		log.Printf("err : %v", err)
		return sqlc.User{}, err
	}
	return user, nil
}

func (ur *userRepository) UpdateUser(context context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.DB.UpdateUser(context, arg)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (ur *userRepository) GetByUuid(ctx context.Context, uuid pgtype.UUID) (sqlc.User, error) {
	user, err := ur.DB.GetByUuid(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (ur *userRepository) ActiveUser(ctx context.Context, userUuid pgtype.UUID) error {
	return ur.DB.ActiveUser(ctx, userUuid)
}
