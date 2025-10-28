package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/internal/dto"
	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/models"
)

type IWebhookService interface {
	ReceiveInvoice(ctx context.Context, req *dto.XenditInvoiceRequest) error
}
type webhookService struct {
	orderRepository repositories.IOrderRepository
}

func (ws *webhookService) ReceiveInvoice(ctx context.Context, req *dto.XenditInvoiceRequest) error {
	orderID, err := strconv.ParseUint(req.ExternalID, 10, 64)
	if err != nil {
		return errors.New("invalid external ID format")
	}

	orderEntity, err := ws.orderRepository.GetOrderByID(ctx, uint(orderID))
	if err != nil {
		return err
	}
	if orderEntity == nil {
		return errors.New("order not found")
	}

	now := time.Now()
	updatedBy := "System"
	orderEntity.OrderStatusCode = models.OrderStatusCodePaid
	orderEntity.UpdatedAt = &now
	orderEntity.UpdatedBy = &updatedBy
	orderEntity.XenditPaidAt = &now
	orderEntity.XenditPaymentChannel = req.PaymentChannel
	orderEntity.XenditPaymentMethod = req.PaymentMethod

	err = ws.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return err
	}

	return nil
}

func NewWebhookService(orderRepository repositories.IOrderRepository) IWebhookService {
	return &webhookService{
		orderRepository: orderRepository,
	}
}
