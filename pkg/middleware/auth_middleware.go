package middleware

import (
	"context"
	"log"

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

	log.Printf("Accessing endpoint: %s", info.FullMethod)


	if am.isPublicEndpoint(info.FullMethod) {
		log.Println("Public endpoint, skipping auth")
		return handler(ctx, req)
	}


	tokenStr, err := utils.ExtractTokenFromContext(ctx)
	if err != nil {
		log.Printf("Failed to extract token: %v", err)
		return nil, status.Error(codes.Unauthenticated, "Invalid or missing authentication token")
	}


	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}


	if utils.IsTokenExpired(claims) {
		log.Println("Token has expired")
		return nil, status.Error(codes.Unauthenticated, "Token has expired")
	}


	if _, found := am.cacheService.Get(tokenStr); found {
		log.Println("Token has been blacklisted")
		return nil, status.Error(codes.Unauthenticated, "Token has been invalidated")
	}


	ctx = utils.InjectClaimsToContext(ctx, claims)
	log.Printf("Authenticated user: %s (%s)", claims.Email, claims.RoleCode)


	return handler(ctx, req)
}

func (am *authMiddleware) isPublicEndpoint(method string) bool {
	publicEndpoints := []string{
		"/auth.AuthService/Register",
		"/auth.AuthService/Login",
	}

	for _, endpoint := range publicEndpoints {
		if method == endpoint {
			return true
		}
	}

	return false
}
