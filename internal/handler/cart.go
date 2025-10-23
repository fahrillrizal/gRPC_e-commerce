package handler

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/cart"
)

type CartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService services.ICartService
}

func (ch *CartHandler) AddToCart(ctx context.Context, req *cart.AddToCartRequest) (*cart.AddToCartResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.AddToCartResponse{
			BaseResponse: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.AddToCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *CartHandler) ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.ListCartResponse{
			BaseResponse: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *CartHandler) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.DeleteCartResponse{
			BaseResponse: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.DeleteCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *CartHandler) UpdateCartQty(ctx context.Context, req *cart.UpdateCartQtyRequest) (*cart.UpdateCartQtyResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.UpdateCartQtyResponse{
			BaseResponse: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.UpdateCartQty(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewCartHandler(cartService services.ICartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}
