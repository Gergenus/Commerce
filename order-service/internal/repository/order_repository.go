package repository

import (
	"context"
	"fmt"

	"github.com/Gergenus/commerce/order-service/internal/models"
	"github.com/Gergenus/commerce/order-service/pkg/db"
	"github.com/google/uuid"
)

const PendingStatus = "pending"

type OrderRepository struct {
	db db.PostgresDB
}

func NewOrderRepository(db db.PostgresDB) OrderRepository {
	return OrderRepository{db: db}
}

func (o *OrderRepository) CreateOrder(ctx context.Context, userId uuid.UUID, price float64) (int, error) {
	const op = "repository.CreateOrder"
	var id int
	tx, err := o.db.DB.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	err = tx.QueryRow(ctx, "INSET INTO orders (customer_id, status, price) VALUES($1, $2, $3) RETURNING id", userId.String(), PendingStatus, price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (o *OrderRepository) FillOrder(ctx context.Context, orderId int, products []models.OrderProduct) error {
	const op = "repository.FillOrder"
	for _, product := range products {
		_, err := o.db.DB.Exec(ctx, "INSERT INTO order_goods (order_id, product_id, seller_id, quantity) VALUES($1, $2, $3, $4)", orderId, product.ID, product.SellerID, product.Quantity)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
