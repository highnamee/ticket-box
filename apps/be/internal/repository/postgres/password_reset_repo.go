package postgres

import (
	"context"
	"errors"
	"time"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, reset *domain.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

func (r *PasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	var reset domain.PasswordReset
	err := r.db.WithContext(ctx).First(&reset, "token_hash = ?", tokenHash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrResetTokenNotFound
		}
		return nil, err
	}
	return &reset, nil
}

func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&domain.PasswordReset{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrResetTokenUsed
	}
	return nil
}
