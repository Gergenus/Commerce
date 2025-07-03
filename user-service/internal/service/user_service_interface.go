package service

import (
	"context"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	AddUser(ctx context.Context, user models.User) (*uuid.UUID, error)
	Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error)
}
