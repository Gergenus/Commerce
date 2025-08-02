package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/order-service/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo repository.OrderRepoInterface
	log  *slog.Logger
}

func NewOrderService(repo repository.OrderRepoInterface, log *slog.Logger) *OrderService {
	return &OrderService{
		repo: repo,
		log:  log,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, userId uuid.UUID) (int, error) {
	const op = "service.CreateOrder"
	log := o.log.With(slog.String("op", op))
	// gRPC req to cart service

	// gRPC req to product service

	orderId, err := o.repo.CreateOrder(ctx, userId)
	if err != nil {
		log.Error("failed to create order", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	err = o.repo.FillOrder(ctx)
	if err != nil {
		log.Error("failed to fill order", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
}
