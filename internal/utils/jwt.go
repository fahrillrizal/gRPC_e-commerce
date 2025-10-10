package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/internal/entity"
	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	jwtClaimsKey contextKey = "jwt_claims"
)

var SecretKey []byte

func InitJWTSecret() error {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return errors.New("JWT_SECRET_KEY environment variable is not set")
	}
	SecretKey = []byte(secret)
	return nil
}

func GenerateJWT(user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	if user.Role == nil {
		return "", errors.New("user role cannot be nil")
	}

	claims := entity.JwtClaims{
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		RoleCode: user.Role.Code,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-grpc",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*entity.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*entity.JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func ExtractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata in context")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", errors.New("missing authorization header")
	}

	token, err := ParseBearerToken(authHeaders[0])
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format, expected 'Bearer <token>'")
	}

	if parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization scheme '%s', must be 'Bearer'", parts[0])
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty token in authorization header")
	}

	return token, nil
}

func ExtractAndValidateToken(ctx context.Context) (*entity.JwtClaims, string, error) {
	token, err := ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, "", err
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		return nil, "", err
	}

	return claims, token, nil
}

func GetTokenExpiry(claims *entity.JwtClaims) time.Duration {
	expiryTime := claims.ExpiresAt.Time
	remaining := time.Until(expiryTime)

	if remaining < 0 {
		return 0
	}

	return remaining
}

func IsTokenExpired(claims *entity.JwtClaims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

func GetTokenAge(claims *entity.JwtClaims) time.Duration {
	return time.Since(claims.IssuedAt.Time)
}

func NeedsRefresh(claims *entity.JwtClaims, threshold time.Duration) bool {
	remaining := GetTokenExpiry(claims)
	return remaining > 0 && remaining < threshold
}

func InjectClaimsToContext(ctx context.Context, claims *entity.JwtClaims) context.Context {
	return context.WithValue(ctx, jwtClaimsKey, claims)
}

func GetClaimsFromContext(ctx context.Context) (*entity.JwtClaims, error) {
	claims, ok := ctx.Value(jwtClaimsKey).(*entity.JwtClaims)
	if !ok {
		return nil, errors.New("claims not found in context")
	}
	return claims, nil
}

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func GetUserEmailFromContext(ctx context.Context) (string, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return "", err
	}
	return claims.RoleCode, nil
}

func GetFullNameFromContext(ctx context.Context) (string, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return "", err
	}
	return claims.FullName, nil
}
