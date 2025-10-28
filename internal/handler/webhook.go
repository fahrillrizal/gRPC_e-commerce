package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fahrillrizal/ecommerce-grpc/internal/dto"
	"github.com/fahrillrizal/ecommerce-grpc/internal/services"
	"github.com/gofiber/fiber/v2"
)

type webhookHandler struct {
	webhookService services.IWebhookService
}

func (wh *webhookHandler) ReceiveInvoice(c *fiber.Ctx) error {
	fmt.Println(string(c.Body()))
	var req dto.XenditInvoiceRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}

	err = wh.webhookService.ReceiveInvoice(c.UserContext(), &req)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func NewWebhookHandler(webhookService services.IWebhookService) *webhookHandler {
	return &webhookHandler{
		webhookService: webhookService,
	}
}
