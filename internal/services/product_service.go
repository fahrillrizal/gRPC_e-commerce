package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/common"
	"github.com/fahrillrizal/ecommerce-grpc/pb/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error)
	DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error)
	UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error)
	DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
	ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error)
	HighlightProducts(ctx context.Context, req *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error)
}

type productService struct {
	productRepository repositories.IProductRepository
	cloudinaryUtils   utils.ICloudinaryUtils
}

func (ps *productService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	if claims.RoleCode != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "only admin can create product")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name is required")
	}

	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product price must be greater than 0")
	}

	imageURL := req.ImageUrl

	if len(req.ImageData) > 0 {
		reader := bytes.NewReader(req.ImageData)

		filename := req.ImageFilename
		if filename == "" {
			filename = req.Name
		}

		uploadedURL, err := ps.cloudinaryUtils.UploadImage(ctx, reader, filename)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to upload image: %v", err))
		}
		imageURL = uploadedURL
	}

	if imageURL == "" {
		imageURL = "https://via.placeholder.com/400x400?text=No+Image"
	}

	newProduct := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    imageURL,
	}

	newProduct.CreatedBy = claims.FullName

	err = ps.productRepository.CreateProduct(ctx, newProduct)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create product: %v", err))
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Product created successfully"),
		Id:   fmt.Sprintf("%d", newProduct.ID),
	}, nil
}

func (ps *productService) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	res, err := ps.productRepository.GetProductByID(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &product.DetailProductResponse{
		Base:        utils.SuccessResponse("Product retrieved successfully"),
		Id:          uint64(res.ID),
		Name:        res.Name,
		Description: res.Description,
		Price:       res.Price,
		ImageUrl:    res.ImageURL,
	}, nil
}

func (ps *productService) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	if claims.RoleCode != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "only admin can update product")
	}

	productID := uint(req.Id)
	existingProduct, err := ps.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	if req.Name != "" {
		if len(req.Name) > 255 {
			return nil, status.Error(codes.InvalidArgument, "name must not exceed 255 characters")
		}
		existingProduct.Name = req.Name
	}

	if req.Description != "" {
		if len(req.Description) > 1000 {
			return nil, status.Error(codes.InvalidArgument, "description must not exceed 1000 characters")
		}
		existingProduct.Description = req.Description
	}

	if req.Price > 0 {
		existingProduct.Price = req.Price
	}

	imageURL := req.ImageUrl
	oldImageURL := existingProduct.ImageURL

	if len(req.ImageData) > 0 {
		reader := bytes.NewReader(req.ImageData)

		filename := req.ImageFilename
		if filename == "" {
			filename = existingProduct.Name
		}

		uploadedURL, err := ps.cloudinaryUtils.UploadImage(ctx, reader, filename)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to upload new image: %v", err))
		}
		imageURL = uploadedURL

		if oldImageURL != "" && oldImageURL != imageURL {
			oldPublicID := ps.cloudinaryUtils.ExtractPublicIDFromURL(oldImageURL)
			if oldPublicID != "" {

				go func() {
					deleteErr := ps.cloudinaryUtils.DeleteImage(context.Background(), oldPublicID)
					if deleteErr != nil {

						fmt.Printf("Warning: failed to delete old image %s: %v\n", oldPublicID, deleteErr)
					}
				}()
			}
		}
	}

	if imageURL != "" {
		existingProduct.ImageURL = imageURL
	}

	existingProduct.UpdatedBy = &claims.FullName

	err = ps.productRepository.UpdateProduct(ctx, existingProduct)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update product: %v", err))
	}

	return &product.UpdateProductResponse{
		Base:        utils.SuccessResponse("Product updated successfully"),
		Id:          uint64(existingProduct.ID),
		Name:        existingProduct.Name,
		Description: existingProduct.Description,
		Price:       existingProduct.Price,
		ImageUrl:    existingProduct.ImageURL,
	}, nil
}

func (ps *productService) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	if claims.RoleCode != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "only admin can delete product")
	}

	productID := uint(req.Id)
	existingProduct, err := ps.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	existingProduct.DeletedBy = &claims.FullName
	existingProduct.UpdatedBy = &claims.FullName

	err = ps.productRepository.DeleteProduct(ctx, existingProduct)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete product: %v", err))
	}

	if existingProduct.ImageURL != "" {
		publicID := ps.cloudinaryUtils.ExtractPublicIDFromURL(existingProduct.ImageURL)
		if publicID != "" {

			go func() {
				deleteErr := ps.cloudinaryUtils.DeleteImage(context.Background(), publicID)
				if deleteErr != nil {

					fmt.Printf("Warning: failed to delete image %s from Cloudinary: %v\n", publicID, deleteErr)
				} else {
					fmt.Printf("Successfully deleted image %s from Cloudinary\n", publicID)
				}
			}()
		}
	}

	return &product.DeleteProductResponse{
		Base: utils.SuccessResponse("Product deleted successfully"),
	}, nil
}

func (ps *productService) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	if req.Pagination == nil {
		req.Pagination = &common.PaginationRequest{
			CurrentPage: 1,
			PerPage:     10,
		}
	}

	products, pagination, err := ps.productRepository.GetProductsPagination(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get products: %v", err))
	}

	var productItems []*product.ListProductResponseItem
	for _, p := range products {
		productItems = append(productItems, &product.ListProductResponseItem{
			Id:          uint64(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ImageUrl:    p.ImageURL,
		})
	}

	return &product.ListProductResponse{
		Base:       utils.SuccessResponse("Products retrieved successfully"),
		Pagination: pagination,
		Data:       productItems,
	}, nil
}

func (ps *productService) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	if claims.RoleCode != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "only admin can access this resource")
	}

	if req.Pagination == nil {
		req.Pagination = &common.PaginationRequest{
			CurrentPage: 1,
			PerPage:     10,
		}
	}

	products, pagination, err := ps.productRepository.GetProductsPaginationAdmin(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get products: %v", err))
	}

	var productItems []*product.ListProductAdminResponseItem
	for _, p := range products {
		item := &product.ListProductAdminResponseItem{
			Id:          uint64(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ImageUrl:    p.ImageURL,
		}

		productItems = append(productItems, item)
	}

	return &product.ListProductAdminResponse{
		Base:       utils.SuccessResponse("Products retrieved successfully"),
		Pagination: pagination,
		Data:       productItems,
	}, nil
}

func (ps *productService) HighlightProducts(ctx context.Context, req *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error) {
	products, err := ps.productRepository.GetHighlightedProducts(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get highlighted products: %v", err))
	}

	var productItems []*product.HighlightProductsResponseItem
	for _, p := range products {
		productItems = append(productItems, &product.HighlightProductsResponseItem{
			Id:          uint64(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			ImageUrl:    p.ImageURL,
		})
	}

	return &product.HighlightProductsResponse{
		Base: utils.SuccessResponse("Highlighted products retrieved successfully"),
		Data: productItems,
	}, nil
}

func NewProductService(
	productRepository repositories.IProductRepository,
	cloudinaryUtils utils.ICloudinaryUtils,
) IProductService {
	return &productService{
		productRepository: productRepository,
		cloudinaryUtils:   cloudinaryUtils,
	}
}
