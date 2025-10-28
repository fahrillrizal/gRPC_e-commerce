package services

import (
	"context"
	"fmt"
	stdos "os"
	"strconv"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/order"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
	ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error)
	ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error)
	DetailOrder(ctx context.Context, req *order.DetailOrderRequest) (*order.DetailOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error)
}

type orderService struct {
	orderRepository   repositories.IOrderRepository
	productRepository repositories.IProductRepository
}

func (os *orderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	tx, err := os.orderRepository.BeginTransaction(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txOrderRepo := os.orderRepository.WithTx(tx)

	numbering, err := txOrderRepo.GetNumbering(ctx, "order")
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "failed to get order numbering")
	}

	var productIds = make([]string, len(req.Products))
	for i := range req.Products {
		productIds[i] = fmt.Sprint(req.Products[i].ProductId)
	}

	products, err := os.productRepository.GetProductsByIDs(ctx, productIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	productMap := make(map[uint64]*models.Product)
	for i := range products {
		productMap[uint64(products[i].ID)] = products[i]
	}

	var total float64 = 0
	for _, p := range req.Products {
		product, exists := productMap[p.ProductId]
		if !exists {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "product with id %d not found", p.ProductId)
		}
		total += product.Price * float64(p.Quantity)
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	orderEntity := models.Order{
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
		UserID:          claims.UserID,
		OrderStatusCode: models.OrderStatusCodeUnpaid,
		UserFullName:    req.FullName,
		Address:         req.Address,
		PhoneNumber:     req.PhoneNumber,
		Notes:           req.Notes,
		Total:           total,
		ExpiredAt:       &expiredAt,
		BaseModel: models.BaseModel{
			CreatedAt: now,
			CreatedBy: claims.FullName,
		},
	}

	err = txOrderRepo.CreateOrder(ctx, &orderEntity)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	invoiceItems := make([]xendit.InvoiceItem, 0)
	for _, p := range req.Products {
		prod := productMap[p.ProductId]
		if prod != nil {
			invoiceItems = append(invoiceItems, xendit.InvoiceItem{
				Name:     prod.Name,
				Price:    prod.Price,
				Quantity: int(p.Quantity),
			})
		}
	}

	frontendURL := stdos.Getenv("FRONTEND_URL")

	xenditInvoice, xenditErr := invoice.CreateWithContext(ctx, &invoice.CreateParams{
		ExternalID: fmt.Sprint(orderEntity.ID),
		Amount:     total,
		Customer: xendit.InvoiceCustomer{
			GivenNames: claims.FullName,
		},
		Currency:           "IDR",
		SuccessRedirectURL: fmt.Sprintf("%s/checkout/%d/success", frontendURL, orderEntity.ID),
		Items:              invoiceItems,
	})
	if xenditErr != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "failed to create xendit invoice: %v", xenditErr)
	}

	if xenditInvoice.ID != "" {
		orderEntity.XenditInvoiceID = xenditInvoice.ID
	}
	if xenditInvoice.InvoiceURL != "" {
		orderEntity.XenditInvoiceUrl = xenditInvoice.InvoiceURL
	}

	err = txOrderRepo.UpdateOrder(ctx, &orderEntity)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, p := range req.Products {
		product := productMap[p.ProductId]
		subtotal := product.Price * float64(p.Quantity)

		var orderItem = models.OrderItem{
			ProductID:    uint(p.ProductId),
			ProductName:  product.Name,
			ProductImage: product.ImageURL,
			ProductPrice: product.Price,
			Quantity:     int(p.Quantity),
			Subtotal:     subtotal,
			OrderID:      orderEntity.ID,
			BaseModel: models.BaseModel{
				CreatedAt: now,
				CreatedBy: claims.FullName,
			},
		}

		err = txOrderRepo.CreateOrderItem(ctx, &orderItem)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	numbering.Number++
	err = txOrderRepo.UpdateNumbering(ctx, numbering)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "failed to update order numbering")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to commit transaction")
	}

	return &order.CreateOrderResponse{
		Base:    utils.SuccessResponse("Create order success"),
		OrderId: fmt.Sprint(orderEntity.ID),
	}, nil
}

func (os *orderService) ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	if claims.RoleCode != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "only admin can delete product")
	}

	orders, metadata, err := os.orderRepository.GetListOrderAdmin(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderAdminResponseItem, 0)
	for _, o := range orders {

		products := make([]*order.ListOrderAdminResponseItemProduct, 0)
		for _, oi := range o.Items {
			products = append(products, &order.ListOrderAdminResponseItemProduct{
				Id:       uint64(oi.ProductID),
				Name:     oi.ProductName,
				Price:    oi.ProductPrice,
				Quantity: int64(oi.Quantity),
			})
		}

		items = append(items, &order.ListOrderAdminResponseItem{
			Id:         fmt.Sprint(o.ID),
			Number:     o.Number,
			Customer:   o.UserFullName,
			StatusCode: o.OrderStatusCode,
			Total:      o.Total,
			CreatedAt:  nil,
			Products:   products,
		})
	}

	return &order.ListOrderAdminResponse{
		Base:       utils.SuccessResponse("List order admin success"),
		Pagination: metadata,
		Orders:     items,
	}, nil
}

func (os *orderService) ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	orders, metadata, err := os.orderRepository.GetListOrder(ctx, claims.UserID, req.Pagination)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderResponseItem, 0)
	for _, o := range orders {

		products := make([]*order.ListOrderResponseItemProduct, 0)
		for _, oi := range o.Items {
			products = append(products, &order.ListOrderResponseItemProduct{
				Id:       uint64(oi.ProductID),
				Name:     oi.ProductName,
				Price:    oi.ProductPrice,
				Quantity: int64(oi.Quantity),
			})
		}

		items = append(items, &order.ListOrderResponseItem{
			Id:               fmt.Sprint(o.ID),
			Number:           o.Number,
			Customer:         o.UserFullName,
			StatusCode:       o.OrderStatusCode,
			Total:            o.Total,
			CreatedAt:        nil,
			Products:         products,
			XenditInvoiceUrl: o.XenditInvoiceUrl,
		})
	}

	return &order.ListOrderResponse{
		Base:       utils.SuccessResponse("List order success"),
		Pagination: metadata,
		Orders:     items,
	}, nil
}

func (os *orderService) DetailOrder(ctx context.Context, req *order.DetailOrderRequest) (*order.DetailOrderResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	orderID, err := strconv.ParseUint(req.OrderId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
	}

	orderEntity, err := os.orderRepository.GetOrderByID(ctx, uint(orderID))
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	if claims.RoleCode != "ADMIN" && orderEntity.UserID != claims.UserID {
		return nil, status.Error(codes.PermissionDenied, "you can only view your own orders")
	}

	items := make([]*order.DetailOrderResponseItem, 0)
	for _, oi := range orderEntity.Items {
		items = append(items, &order.DetailOrderResponseItem{
			Id:       uint64(oi.ProductID),
			Name:     oi.ProductName,
			Price:    oi.ProductPrice,
			Quantity: int64(oi.Quantity),
		})
	}

	return &order.DetailOrderResponse{
		Base:             utils.SuccessResponse("Detail order success"),
		Id:               fmt.Sprint(orderEntity.ID),
		Number:           orderEntity.Number,
		UserFullName:     orderEntity.UserFullName,
		Address:          orderEntity.Address,
		PhoneNumber:      orderEntity.PhoneNumber,
		Notes:            0,
		OrderStatusCode:  orderEntity.OrderStatusCode,
		CreatedAt:        nil,
		XenditInvoiceUrl: orderEntity.XenditInvoiceUrl,
		Items:            items,
	}, nil
}

func (os *orderService) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	claims, err := utils.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	orderID, err := strconv.ParseUint(req.OrderId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
	}

	orderEntity, err := os.orderRepository.GetOrderByID(ctx, uint(orderID))
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	isAdmin := claims.RoleCode == "ADMIN"
	isOwner := orderEntity.UserID == claims.UserID

	if !isAdmin && !isOwner {
		return nil, status.Error(codes.PermissionDenied, "you can only update your own orders")
	}

	currentStatus := orderEntity.OrderStatusCode
	newStatus := req.NewStatusCode

	allowedTransitions := map[string]map[string]bool{
		models.OrderStatusCodeUnpaid: {
			models.OrderStatusCodePaid:      true,
			models.OrderStatusCodeCompleted: true,
		},
		models.OrderStatusCodePaid: {
			models.OrderStatusCodeShipped: true,
		},
		models.OrderStatusCodeShipped: {
			models.OrderStatusCodeCompleted: true,
		},
	}

	if transitions, exists := allowedTransitions[currentStatus]; !exists || !transitions[newStatus] {
		return nil, status.Errorf(codes.InvalidArgument, "invalid status transition from %s to %s", currentStatus, newStatus)
	}

	switch {
	case currentStatus == models.OrderStatusCodeUnpaid && newStatus == models.OrderStatusCodePaid:

		if !isAdmin {
			return nil, status.Error(codes.PermissionDenied, "only admin can mark order as paid")
		}

	case currentStatus == models.OrderStatusCodePaid && newStatus == models.OrderStatusCodeShipped:

		if !isAdmin {
			return nil, status.Error(codes.PermissionDenied, "only admin can mark order as shipped")
		}

	case currentStatus == models.OrderStatusCodeUnpaid && newStatus == models.OrderStatusCodeCompleted:

	case currentStatus == models.OrderStatusCodeShipped && newStatus == models.OrderStatusCodeCompleted:

	}

	now := time.Now()
	orderEntity.OrderStatusCode = newStatus
	orderEntity.UpdatedAt = &now
	orderEntity.UpdatedBy = &claims.FullName

	err = os.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update order status")
	}

	return &order.UpdateOrderStatusResponse{
		Base: utils.SuccessResponse("Order status updated successfully"),
	}, nil
}

func NewOrderService(orderRepository repositories.IOrderRepository, productRepository repositories.IProductRepository) IOrderService {
	return &orderService{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}
