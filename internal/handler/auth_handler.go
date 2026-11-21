package handler

import (
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"

	"github.com/gin-gonic/gin"
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

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password to receive access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginParams true "Login credentials"
// @Success      200 {object} dto.LoginResponse "Login successful"
// @Failure      400 {object} utils.AppError "Invalid request body or validation error"
// @Failure      401 {object} utils.AppError "Unauthorized - Invalid credentials"
// @Router       /auth/login [post]
func (au *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginParams
	if err := c.ShouldBindJSON(&input); err != nil {
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

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate a new access token using a valid refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenParams true "Refresh Token payload"
// @Success      200 {object} dto.RefreshTokenResponse "Token refreshed successfully"
// @Failure      400 {object} utils.AppError "Invalid request body"
// @Failure      401 {object} utils.AppError "Unauthorized - Invalid or expired refresh token"
// @Router       /auth/refresh [post]
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

// Logout godoc
// @Summary      User logout
// @Description  Revoke the active refresh token to log the user out
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.RefreshTokenParams true "Refresh Token to revoke"
// @Success      204 {string} string "No Content"
// @Failure      400 {object} utils.AppError "Invalid request body"
// @Failure      401 {object} utils.AppError "Unauthorized"
// @Router       /auth/logout [post]
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
	utils.ResponseSuccess(c, http.StatusNoContent, nil, "Logout Success")
}

// CreateUser godoc
// @Summary      Register a new user
// @Description  Create a new user account and trigger an OTP verification email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserParams true "User registration details"
// @Success      201 {object} dto.CreateUserRes "User created successfully, pending OTP activation"
// @Failure      400 {object} utils.AppError "Validation error or invalid input"
// @Failure      409 {object} utils.AppError "Email already exists"
// @Router       /auth/create [post]
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

// ValidateOTPAndActive godoc
// @Summary      Activate user account via OTP
// @Description  Verify the OTP sent to the user's email to activate the account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.OTPCreateUserParam true "Email and OTP code"
// @Success      204 {string} string "No Content - Account activated successfully"
// @Failure      400 {object} utils.AppError "Invalid OTP or validation error"
// @Router       /auth/activate [post]
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
