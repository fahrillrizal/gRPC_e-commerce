package middleware

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}

func (am *authMiddleware) Middleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if am.isPublicEndpoint(info.FullMethod) {
		return handler(ctx, req)
	}

	tokenStr, err := utils.ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid or missing authentication token")
	}

	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	if utils.IsTokenExpired(claims) {
		return nil, status.Error(codes.Unauthenticated, "Token has expired")
	}

	if _, found := am.cacheService.Get(tokenStr); found {
		return nil, status.Error(codes.Unauthenticated, "Token has been invalidated")
	}

	ctx = utils.InjectClaimsToContext(ctx, claims)

	if am.isAdminOnlyEndpoint(info.FullMethod) && claims.RoleCode != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "Only administrators can perform this action")
	}

	return handler(ctx, req)
}

func (am *authMiddleware) isPublicEndpoint(method string) bool {
	publicEndpoints := []string{
		"/auth.AuthService/Register",
		"/auth.AuthService/Login",
		"/product.ProductService/ListProduct",
		"/product.ProductService/DetailProduct",
		"/product.ProductService/HighlightProducts",
		"/newsletter.NewsletterService/Subscribe",
	}

	for _, endpoint := range publicEndpoints {
		if method == endpoint {
			return true
		}
	}

	return false
}

func (am *authMiddleware) isAdminOnlyEndpoint(method string) bool {
	adminOnlyEndpoints := []string{
		"/product.ProductService/CreateProduct",
		"/product.ProductService/UpdateProduct",
		"/product.ProductService/DeleteProduct",
	}

	for _, endpoint := range adminOnlyEndpoints {
		if method == endpoint {
			return true
		}
	}

	return false
}
