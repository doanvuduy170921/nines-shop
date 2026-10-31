package validation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"nineshop-be/internal/dto"
	"nineshop-be/internal/utils"
	"strconv"
)

func HandlerValidationError(err error) gin.H {
	if validateErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validateErr {
			switch e.Tag() {
			case "required":
				errors[e.Field()] = fmt.Sprintf("%s is required", e.Field())
			case "email":
				errors[e.Field()] = fmt.Sprintf("%s must be a valid email address", e.Field())
			case "min":
				errors[e.Field()] = fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
			case "max":
				errors[e.Field()] = fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param())
			case "strong_pass":
				errors[e.Field()] = fmt.Sprintf("%s must be 8-20 characters long and include lowercase, uppercase, numbers, and special characters", e.Field())
			case "strong_user":
				errors[e.Field()] = fmt.Sprintf("%s can only contain lowercase letters and numbers", e.Field())
			default:
				errors[e.Field()] = fmt.Sprintf("%s is invalid", e.Field())
			}
		}
		return gin.H{
			"code":    400,
			"message": "Invalid input data",
			"errors":  errors,
		}
	}
	return gin.H{
		"code":    400,
		"message": "Invalid request!",
	}
}

func ConvertAndValidateIntParam(ctx *gin.Context, param string, prefix int32) int32 {
	val, err := strconv.Atoi(param)
	if err != nil {
		utils.ResponseError(ctx, err)
		return 0
	}
	if val <= 0 {
		val = dto.Int32ToInt(prefix)
	}
	return int32(val)
}
