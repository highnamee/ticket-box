package validator

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidateNotBlank checks if string field contains at least one non-whitespace character
func ValidateNotBlank(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return strings.TrimSpace(val) != ""
}

// RegisterCustomValidators registers custom validation tags with Gin's validator engine
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("notblank", ValidateNotBlank)
	}
}
