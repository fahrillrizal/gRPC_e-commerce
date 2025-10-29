package handler

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/newsletter"
)

type newsletterHandler struct {
	newsletter.UnimplementedNewsletterServiceServer

	newsletterService services.INewsletterService
}

func (nh *newsletterHandler) Subscribe(ctx context.Context, req *newsletter.SubscribeRequest) (*newsletter.SubscribeResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &newsletter.SubscribeResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := nh.newsletterService.Subscribe(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewNewsletterHandler(newsletterService services.INewsletterService) *newsletterHandler {
	return &newsletterHandler{
		newsletterService: newsletterService,
	}
}
