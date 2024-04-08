package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/crypitor/postmand"
)

// Delivery implements postmand.DeliveryService interface.
type Delivery struct {
	deliveryRepository postmand.DeliveryRepository
}

// Get returns postmand.Delivery by options filter.
func (d Delivery) Get(ctx context.Context, getOptions postmand.RepositoryGetOptions) (*postmand.Delivery, error) {
	return d.deliveryRepository.Get(ctx, getOptions)
}

// List returns a slice of postmand.Delivery by options filter.
func (d Delivery) List(ctx context.Context, listOptions postmand.RepositoryListOptions) ([]*postmand.Delivery, error) {
	return d.deliveryRepository.List(ctx, listOptions)
}

// Create postmand.Delivery on database.
func (d Delivery) Create(ctx context.Context, delivery *postmand.Delivery) error {
	now := time.Now().UTC()
	delivery.ID = uuid.New()
	delivery.ScheduledAt = now
	delivery.Status = postmand.DeliveryStatusPending
	delivery.CreatedAt = now
	delivery.UpdatedAt = now
	return d.deliveryRepository.Create(ctx, delivery)
}

// Update postmand.Delivery on database.
func (d Delivery) Update(ctx context.Context, delivery *postmand.Delivery) error {
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
	storedDelivery, err := d.deliveryRepository.Get(ctx, getOptions)
	if err != nil {
		return err
	}
	delivery.CreatedAt = storedDelivery.CreatedAt
	delivery.UpdatedAt = time.Now().UTC()
	return d.deliveryRepository.Update(ctx, delivery)
}

// Delete postmand.Delivery on database.
func (d Delivery) Delete(ctx context.Context, id postmand.ID) error {
	return d.deliveryRepository.Delete(ctx, id)
}

// NewDelivery will create an implementation of postmand.DeliveryService.
func NewDelivery(deliveryRepository postmand.DeliveryRepository) *Delivery {
	return &Delivery{deliveryRepository: deliveryRepository}
}
