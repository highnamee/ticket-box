package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordReset struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);not null;index" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}

// BeforeCreate hook generates UUID v7 if not set
func (p *PasswordReset) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = id
	}
	return nil
}

// PasswordResetRepository defines contract for password reset token database operations
type PasswordResetRepository interface {
	Create(ctx context.Context, reset *PasswordReset) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordReset, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
}
