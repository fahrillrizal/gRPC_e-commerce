package services

import (
	"context"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/cart"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ICartService interface {
	AddToCart(ctx context.Context, req *cart.AddToCartRequest) (*cart.AddToCartResponse, error)
	ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error)
	UpdateCartQty(ctx context.Context, req *cart.UpdateCartQtyRequest) (*cart.UpdateCartQtyResponse, error)
}

type cartService struct {
	productRepository repositories.IProductRepository
	cartRepository    repositories.ICartRepository
}

func (cs *cartService) AddToCart(ctx context.Context, req *cart.AddToCartRequest) (*cart.AddToCartResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	product, err := cs.productRepository.GetProductByID(ctx, uint(req.ProductId))
	if err != nil {
		return nil, err
	}

	if product == nil {
		return &cart.AddToCartResponse{
			BaseResponse: utils.NotFoundResponse("Product not Found"),
		}, nil
	}

	existingCart, err := cs.cartRepository.GetCartByProductUserID(ctx, uint(req.ProductId), claims.UserID)
	if err != nil {
		return nil, err
	}

	if existingCart != nil {
		now := time.Now()
		existingCart.Quantity += 1
		existingCart.UpdatedAt = &now
		updatedBy := claims.FullName
		existingCart.UpdatedBy = &updatedBy

		err = cs.cartRepository.UpdateCart(ctx, existingCart)
		if err != nil {
			return nil, err
		}

		return &cart.AddToCartResponse{
			BaseResponse: utils.SuccessResponse("Added to cart successfully"),
			Id:           string(rune(existingCart.ID)),
		}, nil
	}

	newCart := &models.Cart{
		ProductID: uint(req.ProductId),
		UserID:    claims.UserID,
		Quantity:  int(req.Quantity),
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			CreatedBy: claims.FullName,
		},
	}

	err = cs.cartRepository.CreateNewCart(ctx, newCart)
	if err != nil {
		return nil, err
	}

	return &cart.AddToCartResponse{
		BaseResponse: utils.SuccessResponse("Added to cart successfully"),
		Id:           string(rune(newCart.ID)),
	}, nil
}

func (cs *cartService) ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	carts, err := cs.cartRepository.GetListCart(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	var cartItems []*cart.ListCartResponseItem = make([]*cart.ListCartResponseItem, 0)
	for _, cartItem := range carts {
		item := &cart.ListCartResponseItem{
			ProductId:       int64(cartItem.ProductID),
			ProductName:     cartItem.Product.Name,
			ProductImageUrl: cartItem.Product.ImageURL,
			ProductPrice:    cartItem.Product.Price,
			Quantity:        int32(cartItem.Quantity),
		}
		cartItems = append(cartItems, item)
	}

	return &cart.ListCartResponse{
		BaseResponse: utils.SuccessResponse("List cart fetched successfully"),
		Items:        cartItems,
	}, nil
}

func (cs *cartService) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	existingCart, err := cs.cartRepository.GetCartById(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	if existingCart == nil {
		return &cart.DeleteCartResponse{
			BaseResponse: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if existingCart.UserID != claims.UserID {
		return &cart.DeleteCartResponse{
			BaseResponse: utils.UnauthorizedResponse("You are not authorized to delete this cart"),
		}, nil
	}

	err = cs.cartRepository.DeleteCart(ctx, req.CartId, claims.FullName)
	if err != nil {
		return nil, err
	}

	return &cart.DeleteCartResponse{
		BaseResponse: utils.SuccessResponse("Cart deleted successfully"),
	}, nil
}

func (cs *cartService) UpdateCartQty(ctx context.Context, req *cart.UpdateCartQtyRequest) (*cart.UpdateCartQtyResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	existingCart, err := cs.cartRepository.GetCartById(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	if existingCart == nil {
		return &cart.UpdateCartQtyResponse{
			BaseResponse: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if existingCart.UserID != claims.UserID {
		return &cart.UpdateCartQtyResponse{
			BaseResponse: utils.UnauthorizedResponse("You are not authorized to update this cart"),
		}, nil
	}

	existingCart.Quantity = int(req.NewQuantity)
	err = cs.cartRepository.UpdateCart(ctx, existingCart)
	if err != nil {
		return nil, err
	}

	return &cart.UpdateCartQtyResponse{
		BaseResponse: utils.SuccessResponse("Cart quantity updated successfully"),
	}, nil
}

func NewCartService(productRepository repositories.IProductRepository, cartRepository repositories.ICartRepository) ICartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
	}
}
