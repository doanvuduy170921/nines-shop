package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
)

type AuthHandler struct {
	service service.AuthService
	user    service.UserService
}

func NewAuthHandler(service service.AuthService, user service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
		user:    user,
	}
}

func (au *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginParams
	if err := c.ShouldBindJSON(&input); err != nil { // Đọc JSON từ request body và map vào struct input
		utils.ResponseError(c, err)
		return
	}
	accessToken, refreshToken, TTL, role, err := au.service.Login(c, input.Email, input.Password)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TTL:          TTL,
		Role:         role,
	}

	utils.ResponseSuccess(c, http.StatusOK, response, "Login Success")
}

func (au *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.RefreshTokenParams
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ResponseError(c, err)
		return

	}
	accessToken, refreshToken, ttl, err := au.service.RefreshToken(c, input.RefreshToken)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	response := dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TTL:          ttl,
	}
	utils.ResponseSuccess(c, http.StatusOK, response, "Refresh Token Success")
}

func (au *AuthHandler) Logout(c *gin.Context) {
	var input dto.RefreshTokenParams
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ResponseError(c, err)
		return
	}
	if err := au.service.Logout(c, input.RefreshToken); err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusNoContent, nil, "Logout ")
}

func (au *AuthHandler) CreateUser(c *gin.Context) {
	var input dto.CreateUserParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, validation.HandlerValidationError(err))
		return

	}
	userParam := dto.MapDtoToParams(input)
	userCreated, err := au.user.CreateUser(c, &userParam)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, userCreated, "Create User Success")
}

func (au *AuthHandler) ValidateOTPAndActive(c *gin.Context) {
	var input dto.OTPCreateUserParam
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, validation.HandlerValidationError(err))
		return
	}
	ctx := c.Request.Context()
	if err := au.user.ActiveUser(ctx, input.Email, input.Otp); err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, http.StatusNoContent, nil, "Active User Success")
}
