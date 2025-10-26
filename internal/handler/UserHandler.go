package handler

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/pagination"
	"nineshop-be/internal/service"
	"nineshop-be/internal/utils"
	"nineshop-be/internal/validation"
	"strconv"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(c *gin.Context) {
	//page := c.DefaultQuery("page", "1")
	//limit := c.DefaultQuery("limit", "10")
	//
	//pageNumber, _ := strconv.Atoi(page)
	//limitNumber, _ := strconv.Atoi(limit)

	users, err := uh.service.GetAllUser(c)
	if err != nil {
		utils.ResponseError(c, err)
		return

	}
	userRes := dto.MapUsersToResponse(users)
	utils.ResponseSuccess(c, http.StatusOK, userRes, "Get All User Success")
}

func (uh *UserHandler) GetAllByFilter(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	role := c.DefaultQuery("role", "")
	isActiveStr := c.Query("is_active")
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	var isActive *bool = nil
	if isActiveStr != "" {
		val, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			utils.ResponseError(c, err)
			return
		}
		isActive = &val
	}

	users, count, err := uh.service.GetAllUserV2(c, search, role, isActive, utils.StringToInt32(page), utils.StringToInt32(limit))
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	usersRes := dto.MapUsersToResponse(users)
	paginationRes := pagination.PaginationRes{
		Data:      usersRes,
		Total:     int32(count),
		Page:      utils.StringToInt32(page),
		Limit:     utils.StringToInt32(limit),
		TotalPage: int32(math.Ceil(float64(count) / utils.StringToFloat64(limit))),
	}

	utils.ResponseSuccess(c, http.StatusOK, paginationRes, "Get All User Success")
}

func (uh *UserHandler) SoftDelete(c *gin.Context) {
	uuidStr := c.Param("uuid")

	uuid, err := utils.StringToPgUuid(uuidStr)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	user, err := uh.service.SoftDeleteUser(c, uuid)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	userRes := dto.MapUserToResponse(user)

	utils.ResponseSuccess(c, http.StatusOK, userRes, "SoftDelete User Success")

}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	uuidStr := c.Param("uuid")

	var input dto.UpdateUserParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, validation.HandlerValidationError(err))
		return
	}
	userParam := dto.MapUpdateUserDtoToUser(uuidStr, input)

	user, err := uh.service.UpdateUser(c, userParam)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	userRes := dto.UpdateUserToResponse(user)

	utils.ResponseSuccess(c, http.StatusOK, userRes, "Update User Success")
}

func (uh *UserHandler) GetByUuid(c *gin.Context) {
	uuidStr := c.Param("uuid")
	uuid, err := utils.StringToPgUuid(uuidStr)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	user, err := uh.service.GetByUuid(c, uuid)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	userRes := dto.MapUserToResponse(user)
	utils.ResponseSuccess(c, http.StatusOK, userRes, "Get User Success")
}
