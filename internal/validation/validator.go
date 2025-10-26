package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 20 {
		return false
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[\W_]`).MatchString(password) {
		return false
	}
	return true
}

func UsernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	re := regexp.MustCompile(`^[a-z0-9]+$`)
	return re.MatchString(username)
}

func RegisterValidation(v *validator.Validate) {
	v.RegisterValidation("strong_pass", PasswordValidator)
	v.RegisterValidation("strong_user", UsernameValidator)

}
