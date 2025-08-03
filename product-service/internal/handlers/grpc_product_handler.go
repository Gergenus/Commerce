package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/service"
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
		if errors.Is(err, service.ErrStockNotFound) {
			return &proto.AvailablilityResponse{Availablility: false}, nil
		}
		return &proto.AvailablilityResponse{}, ErrInternal

	}
	if stock < int(in.GetStock()) {
		return &proto.AvailablilityResponse{Availablility: false}, nil
	}
	return &proto.AvailablilityResponse{Availablility: true}, nil
}

func (p *ProductHandler) ReserveOrder(ctx context.Context, in *proto.ReserveOrderRequest) (*proto.ReserveOrderResponse, error) {
	products := []models.ProductsToReserve{}
	for _, d := range in.OrderProducts {
		products = append(products, models.ProductsToReserve{ID: int(d.ProductId), Stock: int(d.Stock)})
	}
	reservedProducts, err := p.service.ReserveProducts(ctx, products)
	if err != nil {
		return &proto.ReserveOrderResponse{IsReserved: false}, ErrInternal
	}
	var fullPrice float64
	responseProducts := []*proto.ProductSeller{}
	for _, d := range reservedProducts {
		fullPrice += d.Price
		responseProducts = append(responseProducts, &proto.ProductSeller{ProductId: int64(d.ID), Stock: int64(d.Stock), SellerId: d.SellerID.String()})
	}
	return &proto.ReserveOrderResponse{ProductsSeller: responseProducts, Price: float32(fullPrice), IsReserved: true}, nil
}
