package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/Gergenus/commerce/product-service/proto"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
	ErrInternal       = errors.New("internal error")
)

func (p *ProductHandler) IsAvailable(ctx context.Context, in *proto.AvailablilityRequest) (*proto.AvailablilityResponse, error) {
	if in.GetProductId() == "" || in.GetStock() <= 0 {
		return &proto.AvailablilityResponse{}, ErrInvalidPayload
	}
	productId, err := strconv.Atoi(in.GetProductId())
	if err != nil {
		return &proto.AvailablilityResponse{}, ErrInvalidPayload
	}
	stock, err := p.service.GetStockByID(ctx, productId)
	if err != nil {
		return &proto.AvailablilityResponse{}, ErrInternal

	}
	if stock < int(in.GetStock()) {
		return &proto.AvailablilityResponse{Availablility: false}, nil
	}
	return &proto.AvailablilityResponse{Availablility: true}, nil
}
