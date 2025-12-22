package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/utils"
)

type CreateUserParams struct {
	Name     string `json:"name" binding:"required,min=3,max=50,strong_user"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,strong_pass"`
	Phone    string `json:"phone" binding:"omitempty"`
	Address  string `json:"address" binding:"omitempty,max=255"`
	Role     string `json:"role" binding:"omitempty"`
}

type UpdateUserParams struct {
	Name     string `json:"name" binding:"required,min=3,max=50,strong_user"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Address  string `json:"address" binding:"omitempty,max=255"`
	Role     string `json:"role" binding:"omitempty"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

func MapUpdateUserDtoToUser(userUuid string, input UpdateUserParams) sqlc.UpdateUserParams {
	if input.Role == "" {
		input.Role = "customer"
	}
	val, _ := utils.StringToPgUuid(userUuid)

	return sqlc.UpdateUserParams{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Address:  input.Address,
		Role:     &input.Role,
		IsActive: *input.IsActive,
		UserUuid: val,
	}
}

func MapDtoToParams(input CreateUserParams) sqlc.CreateUserParams {
	if input.Role == "" {
		input.Role = "customer"
	}
	active := false
	return sqlc.CreateUserParams{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Phone:    &input.Phone,
		Address:  &input.Address,
		Role:     &input.Role,
		IsActive: &active,
	}
}

func StringToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{
			String: s,
			Valid:  false,
		}
	}
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

func PgTextToString(s pgtype.Text) string {
	if !s.Valid {
		return "" // Trả về chuỗi rỗng cho nil hoặc giá trị NULL
	}
	return s.String
}

type CreateUserResponse struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Created  string `json:"created" binding:"required"`
	Updated  string `json:"updated" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
	UserUuid string `json:"user_uuid" binding:"required"`
}

type UpdateUserResponse struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Role     string `json:"role" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

func MapUserToResponse(input sqlc.User) CreateUserResponse {
	return CreateUserResponse{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    ptrToString(input.Phone),
		Address:  ptrToString(input.Address),
		Role:     ptrToString(input.Role),
		Created:  pgTimeToString(input.CreatedAt),
		Updated:  pgTimeToString(input.UpdatedAt),
		IsActive: *input.IsActive,
		UserUuid: utils.PgTypeUuidToString(input.UserUuid),
	}
}

func UpdateUserToResponse(input sqlc.User) UpdateUserResponse {
	return UpdateUserResponse{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    ptrToString(input.Phone),
		Address:  ptrToString(input.Address),
		Role:     ptrToString(input.Role),
		IsActive: *input.IsActive,
	}
}

func MapUsersToResponse(input []sqlc.User) []CreateUserResponse {
	res := make([]CreateUserResponse, 0, len(input))
	for _, user := range input {
		res = append(res, MapUserToResponse(user))
	}
	return res
}

func ptrToString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func pgTimeToString(input pgtype.Timestamp) string {
	if !input.Valid {
		return ""
	}
	return input.Time.Format("2006-01-02 15:04:05")

}

type CreateUserRes struct {
	Email     string `json:"email" binding:"required"`
	ExpiredAt string `json:"expired_at" binding:"required"`
}
