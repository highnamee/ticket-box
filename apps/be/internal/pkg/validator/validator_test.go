package validator_test

import (
	"testing"

	customValidator "ticket-box-be/internal/pkg/validator"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestPayload struct {
	UUIDField  string `validate:"uuidv7"`
	TitleField string `validate:"notblank"`
}

func TestCustomValidators(t *testing.T) {
	validate := validator.New()
	_ = validate.RegisterValidation("uuidv7", customValidator.ValidateUUIDv7)
	_ = validate.RegisterValidation("notblank", customValidator.ValidateNotBlank)

	t.Run("uuidv7 validator", func(t *testing.T) {
		validV7, err := uuid.NewV7()
		require.NoError(t, err)

		v4 := uuid.New() // UUID v4

		// Valid UUID v7
		payload1 := TestPayload{UUIDField: validV7.String(), TitleField: "Valid Title"}
		assert.NoError(t, validate.Struct(payload1))

		// Invalid UUID v4 (should fail uuidv7 check)
		payload2 := TestPayload{UUIDField: v4.String(), TitleField: "Valid Title"}
		assert.Error(t, validate.Struct(payload2))

		// Invalid string format
		payload3 := TestPayload{UUIDField: "not-a-uuid", TitleField: "Valid Title"}
		assert.Error(t, validate.Struct(payload3))
	})

	t.Run("notblank validator", func(t *testing.T) {
		validV7, _ := uuid.NewV7()

		// Valid non-blank
		payload1 := TestPayload{UUIDField: validV7.String(), TitleField: "Event Name"}
		assert.NoError(t, validate.Struct(payload1))

		// Whitespaces only (should fail)
		payload2 := TestPayload{UUIDField: validV7.String(), TitleField: "   "}
		assert.Error(t, validate.Struct(payload2))

		// Empty string (should fail)
		payload3 := TestPayload{UUIDField: validV7.String(), TitleField: ""}
		assert.Error(t, validate.Struct(payload3))
	})
}
