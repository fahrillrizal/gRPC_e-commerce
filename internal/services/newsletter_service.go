package services

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/internal/repositories"
	"github.com/fahrillrizal/ecommerce-grpc/internal/utils"
	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/newsletter"
)

type INewsletterService interface {
	Subscribe(ctx context.Context, req *newsletter.SubscribeRequest) (*newsletter.SubscribeResponse, error)
}

type newsletterService struct {
	newsletterRepository repositories.INewsletterRepository
}

func (ns *newsletterService) Subscribe(ctx context.Context, req *newsletter.SubscribeRequest) (*newsletter.SubscribeResponse, error) {
	newsletterEntity, err := ns.newsletterRepository.GetNewsletterByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if newsletterEntity != nil {
		return &newsletter.SubscribeResponse{
			Base: utils.SuccessResponse("Subscribe newsletter success"),
		}, nil
	}

	newNewsletterEntity := &models.Newsletter{
		ID:       0,
		FullName: req.Name,
		Email:    req.Email,
	}

	err = ns.newsletterRepository.CreateNewsletter(ctx, newNewsletterEntity)
	if err != nil {
		return nil, err
	}

	return &newsletter.SubscribeResponse{
		Base: utils.SuccessResponse("Subscribe newsletter success"),
	}, nil
}

func NewNewsletterService(newsletterRepository repositories.INewsletterRepository) INewsletterService {
	return &newsletterService{
		newsletterRepository: newsletterRepository,
	}
}
