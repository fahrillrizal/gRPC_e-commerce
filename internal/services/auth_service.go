package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/auth"
	"github.com/fahrillrizal/ecommerce-grpc/pb/common"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}

type authService struct {
	authRepository repositories.IAuthRepository
	cacheService   *gocache.Cache
}

func (as *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	passwordErrors := utils.ValidatePasswordStrength(req.GetPassword())
	if len(passwordErrors) > 0 {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(passwordErrors),
		}, nil
	}

	if matchError := utils.ValidatePasswordMatch(req.GetPassword(), req.GetPasswordConfirmation()); matchError != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse([]*common.ValidationError{matchError}),
		}, nil
	}

	user, err := as.authRepository.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exists"),
		}, nil
	}

	role, err := as.authRepository.GetRoleByCode(ctx, "CUSTOMER")
	if err != nil {
		return nil, err
	}
	if role == nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Role not found"),
		}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 10)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: string(hashedPassword),
		RoleID:   &role.ID,
		BaseModel: models.BaseModel{
			CreatedBy: req.GetFullName(),
		},
	}

	if err := as.authRepository.CreateUser(ctx, newUser); err != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Failed to create user"),
		}, nil
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("Registration successful."),
	}, nil

}

func (as *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := as.authRepository.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("User not found"),
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword())); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Unathenticated")
		}
		return nil, err
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Base:  utils.SuccessResponse("Login successful."),
		Token: token,
	}, nil
}

func (as *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	claims, token, err := utils.ExtractAndValidateToken(ctx)
	if err != nil {
		return &auth.LogoutResponse{
			Base: utils.UnauthorizedResponse("Invalid or missing authentication token"),
		}, nil
	}

	if _, found := as.cacheService.Get(token); found {
		return &auth.LogoutResponse{
			Base: utils.BadRequestResponse("Token already invalidated"),
		}, nil
	}

	ttl := utils.GetTokenExpiry(claims)
	if ttl <= 0 {
		return &auth.LogoutResponse{
			Base: utils.BadRequestResponse("Token already expired"),
		}, nil
	}

	as.cacheService.Set(token, "blacklisted", ttl)

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout successful"),
	}, nil
}

func (as *authService) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	passwordErrors := utils.ValidatePasswordStrength(req.GetNewPassword())
	if len(passwordErrors) > 0 {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse(passwordErrors),
		}, nil
	}

	if matchError := utils.ValidatePasswordMatch(req.GetNewPassword(), req.GetNewPasswordConfirmation()); matchError != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse([]*common.ValidationError{matchError}),
		}, nil
	}

	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.UnauthorizedResponse("Invalid authentication"),
		}, nil
	}

	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.InternalServerErrorResponse("Failed to get user data"),
		}, nil
	}

	if user == nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("User not found"),
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetCurrentPassword())); err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("Current password is incorrect"),
		}, nil
	}

	if req.GetCurrentPassword() == req.GetNewPassword() {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("New password must be different from current password"),
		}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.InternalServerErrorResponse("Failed to process new password"),
		}, nil
	}

	if err := as.authRepository.UpdateUserPassword(ctx, user.ID, string(hashedPassword), claims.FullName); err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.InternalServerErrorResponse("Failed to update password"),
		}, nil
	}

	return &auth.ChangePasswordResponse{
		Base: utils.SuccessResponse("Password changed successfully"),
	}, nil
}

func (as *authService) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return &auth.GetProfileResponse{
			Base: utils.UnauthorizedResponse("Invalid authentication"),
		}, nil
	}

	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return &auth.GetProfileResponse{
			Base: utils.InternalServerErrorResponse("Failed to get user data"),
		}, nil
	}

	if user == nil {
		return &auth.GetProfileResponse{
			Base: utils.BadRequestResponse("User not found"),
		}, nil
	}

	if user.Role == nil {
		return &auth.GetProfileResponse{
			Base: utils.InternalServerErrorResponse("User role not found"),
		}, nil
	}

	return &auth.GetProfileResponse{
		Base:        utils.SuccessResponse("Profile retrieved successfully"),
		Id:          strconv.FormatUint(uint64(user.ID), 10),
		FullName:    user.FullName,
		Email:       user.Email,
		RoleName:    user.Role.Name,
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil

}

func NewAuthService(authRepository repositories.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}
