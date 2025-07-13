package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/Gergenus/commerce/product-service/proto"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
)

func (p *ProductHandler) IsAvailable(ctx context.Context, in *proto.AvailablilityRequest) (*proto.AvailablilityResponse, error) {
	if in.GetProductId() == "" || in.GetStock() <= 0 {
		return &proto.AvailablilityResponse{}, ErrInvalidPayload
	}
	productId, err := strconv.Atoi(in.GetProductId())
	if err != nil {
		return &proto.AvailablilityResponse{}, ErrInvalidPayload
	}
	p.service.GetStockByID(ctx, productId)
}
