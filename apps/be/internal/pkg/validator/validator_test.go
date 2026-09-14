package validator_test

import (
	"testing"

	customValidator "ticket-box-be/internal/pkg/validator"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	TitleField string `validate:"notblank"`
}

func TestCustomValidators(t *testing.T) {
	validate := validator.New()
	_ = validate.RegisterValidation("notblank", customValidator.ValidateNotBlank)

	t.Run("notblank validator", func(t *testing.T) {
		// Valid non-blank
		payload1 := TestPayload{TitleField: "Event Name"}
		assert.NoError(t, validate.Struct(payload1))

		// Whitespaces only (should fail)
		payload2 := TestPayload{TitleField: "   "}
		assert.Error(t, validate.Struct(payload2))

		// Empty string (should fail)
		payload3 := TestPayload{TitleField: ""}
		assert.Error(t, validate.Struct(payload3))
	})
}
