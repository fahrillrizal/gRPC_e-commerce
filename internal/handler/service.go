package handler

import (
	"context"
	"fmt"

	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloworldServiceServer
}

func (sh *serviceHandler) SayHello(ctx context.Context, req *service.HelloRequest) (*service.HelloResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &service.HelloResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	return &service.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
		Base:    utils.SuccessResponse("Success"),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}
