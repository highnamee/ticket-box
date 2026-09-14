package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"ticket-box-be/internal/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	t.Run("development logger", func(t *testing.T) {
		l := logger.InitLogger("development")
		require.NotNil(t, l)
		assert.Equal(t, slog.Default(), l)
	})

	t.Run("production logger", func(t *testing.T) {
		l := logger.InitLogger("production")
		require.NotNil(t, l)
	})

	t.Run("test logger", func(t *testing.T) {
		l := logger.InitLogger("test")
		require.NotNil(t, l)
	})
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()

	t.Run("request ID context", func(t *testing.T) {
		assert.Equal(t, "", logger.GetRequestID(ctx))

		ctxWithReq := logger.WithRequestID(ctx, "req-12345")
		assert.Equal(t, "req-12345", logger.GetRequestID(ctxWithReq))

		logFromCtx := logger.FromContext(ctxWithReq)
		require.NotNil(t, logFromCtx)
	})

	t.Run("custom logger context", func(t *testing.T) {
		customLogger := slog.Default().With("custom", "val")
		ctxWithLogger := logger.WithLogger(ctx, customLogger)
		assert.Equal(t, customLogger, logger.FromContext(ctxWithLogger))
	})
}
