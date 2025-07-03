package repository

import (
	"context"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	AddUser(ctx context.Context, user models.User) (*uuid.UUID, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	CreateJWTSession(ctx context.Context, user models.User, refreshTokenHash, fingerprint, ip string, expiresIn int64) error
}
