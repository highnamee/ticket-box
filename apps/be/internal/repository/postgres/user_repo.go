package postgres

import (
	"context"
	"errors"
	"strings"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// isUniqueViolation reports whether err is a PostgreSQL unique constraint violation (code 23505).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "LOWER(email) = LOWER(?)", strings.TrimSpace(email)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return domain.ErrEmailAlreadyExists
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
