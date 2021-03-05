package service

import (
	"context"

	"github.com/allisson/postmand"
)

// DeliveryAttempt implements postmand.DeliveryAttemptService interface.
type DeliveryAttempt struct {
	deliveryAttemptRepository postmand.DeliveryAttemptRepository
}

// Get returns postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) Get(ctx context.Context, getOptions postmand.RepositoryGetOptions) (*postmand.DeliveryAttempt, error) {
	return d.deliveryAttemptRepository.Get(ctx, getOptions)
}

// List returns a slice of postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) List(ctx context.Context, listOptions postmand.RepositoryListOptions) ([]*postmand.DeliveryAttempt, error) {
	return d.deliveryAttemptRepository.List(ctx, listOptions)
}

// NewDeliveryAttempt will create an implementation of postmand.DeliveryAttemptService.
func NewDeliveryAttempt(deliveryAttemptRepository postmand.DeliveryAttemptRepository) *DeliveryAttempt {
	return &DeliveryAttempt{deliveryAttemptRepository: deliveryAttemptRepository}
}
