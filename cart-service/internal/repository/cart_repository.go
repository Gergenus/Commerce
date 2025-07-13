package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rds     *redis.Client
	cartTTL time.Duration
}

func NewRedisRepository(rds *redis.Client, cartTTL time.Duration) *RedisRepository {
	return &RedisRepository{rds: rds, cartTTL: cartTTL}
}

func (r RedisRepository) AddToCart(ctx context.Context, UUID, productID string, stock int) error {
	const op = "repository.AddToCart"
	err := r.rds.HSet(ctx, UUID, productID, stock).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = r.rds.Expire(ctx, UUID, r.cartTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r RedisRepository) DeleteFromCart(ctx context.Context, UUID, productID string) error {
	const op = "repository.DeleteFromCart"
	err := r.rds.HDel(ctx, UUID, productID).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r RedisRepository) UpdateStock(ctx context.Context, UUID, productID string, stock int) error {
	const op = "repository.UpdateStock"
	err := r.rds.HSet(ctx, UUID, productID, stock).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
