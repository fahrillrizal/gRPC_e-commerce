package handler

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/order"
)

type orderHandler struct {
	order.UnimplementedOrderServiceServer

	orderService services.IOrderService
}

func (oh *orderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.CreateOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (oh *orderHandler) ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.ListOrderAdminResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.ListOrderAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (oh *orderHandler) ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.ListOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.ListOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (oh *orderHandler) DetailOrder(ctx context.Context, req *order.DetailOrderRequest) (*order.DetailOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.DetailOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.DetailOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (oh *orderHandler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &order.UpdateOrderStatusResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.UpdateOrderStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewOrderHandler(orderService services.IOrderService) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}
