package repository

import (
	"context"

	"github.com/Gergenus/commerce/order-service/internal/models"
	"github.com/google/uuid"
)

type OrderRepoInterface interface {
	CreateOrder(ctx context.Context, userId uuid.UUID, price float64) (int, error)
	FillOrder(ctx context.Context, orderId int, products []models.OrderProduct) error
}
