package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/order-service/internal/models"
	"github.com/Gergenus/commerce/order-service/internal/repository"
	"github.com/Gergenus/commerce/order-service/proto"
	"github.com/google/uuid"
)

var (
	ErrNoCartFound      = errors.New("no cart found")
	ErrOrderNotReserved = errors.New("order not reserved")
)

type OrderService struct {
	repo          repository.OrderRepoInterface
	log           *slog.Logger
	cartClient    proto.OrderServiceClient
	productClient proto.OrderServiceClient
}

func NewOrderService(repo repository.OrderRepoInterface, log *slog.Logger, cartClient proto.OrderServiceClient, productClient proto.OrderServiceClient) *OrderService {
	return &OrderService{
		repo:          repo,
		log:           log,
		cartClient:    cartClient,
		productClient: productClient,
	}
}

// returns orderId
func (o *OrderService) CreateOrder(ctx context.Context, userId uuid.UUID) (int, error) {
	const op = "service.CreateOrder"
	log := o.log.With(slog.String("op", op))
	// gRPC req to cart service
	cartResp, err := o.cartClient.GetCart(ctx, &proto.GetCartRequest{UserId: userId.String()})
	if err != nil {
		log.Error("failed to call gRPC req to cart-service", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !cartResp.Availablility {
		log.Warn("no cart was found")
		return 0, fmt.Errorf("%s: %w", op, ErrNoCartFound)
	}

	products := cartResp.GetOrderProducts()
	// gRPC req to product-service, it returns a slice of products with its seller_id
	ReservedResponse, err := o.productClient.ReserveOrder(ctx, &proto.ReserveOrderRequest{OrderProducts: products})
	if !ReservedResponse.IsReserved {
		return 0, fmt.Errorf("%s: %w", op, ErrOrderNotReserved)
	}
	if err != nil {
		log.Error("failed to call gRPC req to product-service", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	orderId, err := o.repo.CreateOrder(ctx, userId, float64(ReservedResponse.Price))
	if err != nil {
		log.Error("failed to create order", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	convertedProducts := convertProducts(products)

	err = o.repo.FillOrder(ctx, orderId, convertedProducts)
	if err != nil {
		log.Error("failed to fill order", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return orderId, nil
}

func convertProducts(oldProducts []*proto.OrderProduct) []models.OrderProduct {
	products := []models.OrderProduct{}
	for _, d := range oldProducts {
		product := models.OrderProduct{
			ID:    int(d.GetProductId()),
			Stock: int(d.GetStock()),
		}
		products = append(products, product)
	}
	return products
}
