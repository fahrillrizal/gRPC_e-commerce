package handler

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer

	authService services.IAuthService
}

func (sh *authHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &auth.LogoutResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := sh.authService.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {	
	res, err := sh.authService.GetProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewAuthHandler(authService services.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
