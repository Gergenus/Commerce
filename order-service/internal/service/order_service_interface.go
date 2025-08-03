package service

import (
	"context"

	"github.com/google/uuid"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, userId uuid.UUID) (int, error)
}
