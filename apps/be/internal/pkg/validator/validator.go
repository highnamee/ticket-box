package validator

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ValidateUUIDv7 checks if the field is a valid UUIDv7
func ValidateUUIDv7(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val == "" {
		return true // Allow optional fields unless marked with 'required'
	}
	id, err := uuid.Parse(val)
	if err != nil {
		return false
	}
	return id.Version() == 7
}

// ValidateNotBlank checks if string field contains at least one non-whitespace character
func ValidateNotBlank(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return strings.TrimSpace(val) != ""
}

// RegisterCustomValidators registers custom validation tags with Gin's validator engine
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("uuidv7", ValidateUUIDv7)
		_ = v.RegisterValidation("notblank", ValidateNotBlank)
	}
}
