package handler

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/pb/product"
)

type productHandler struct {
	product.UnimplementedProductServiceServer

	productService services.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.DetailProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.DetailProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.UpdateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.UpdateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.DeleteProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.DeleteProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.ListProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.ListProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.ListProductAdminResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.ListProductAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) HighlightProducts(ctx context.Context, req *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &product.HighlightProductsResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.HighlightProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewProductHandler(productService services.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
}
