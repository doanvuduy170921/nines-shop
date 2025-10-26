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
				errors[e.Field()] = fmt.Sprintf("%s là bắt buộc", e.Field())
			case "email":
				errors[e.Field()] = fmt.Sprintf("%s phải đúng định dạng email", e.Field())
			case "min":
				errors[e.Field()] = fmt.Sprintf("%s phải có ít nhất %s ký tự", e.Field(), e.Param())
			case "max":
				errors[e.Field()] = fmt.Sprintf("%s không được vượt quá %s ký tự", e.Field(), e.Param())
			case "strong_pass":
				errors[e.Field()] = fmt.Sprintf("%s phải có độ dài 8-20 ký tự, chứa chữ thường, chữ hoa, số và ký tự đặc biệt", e.Field())
			case "strong_user":
				errors[e.Field()] = fmt.Sprintf("%s chỉ được chứa chữ thường, số", e.Field())
			default:
				errors[e.Field()] = fmt.Sprintf("%s không hợp lệ", e.Field())
			}
		}
		return gin.H{"errors": errors}
	}
	return gin.H{"error": "Yêu cầu không hợp lệ!"}
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
