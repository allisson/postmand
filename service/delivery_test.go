package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/crypitor/postmand"
	"github.com/crypitor/postmand/mocks"
)

func TestDelivery(t *testing.T) {
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		deliveryService := NewDelivery(deliveryRepository)
		expectedDelivery := &postmand.Delivery{ID: uuid.New()}
		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": expectedDelivery.ID}}

		deliveryRepository.On("Get", mock.Anything, getOptions).Return(expectedDelivery, nil)
		delivery, err := deliveryService.Get(ctx, getOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDelivery, delivery)
		deliveryRepository.AssertExpectations(t)
	})

	t.Run("List", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		webhookService := NewDelivery(deliveryRepository)
		expectedDelivery := &postmand.Delivery{ID: uuid.New()}
		listOptions := postmand.RepositoryListOptions{Filters: map[string]interface{}{"id": expectedDelivery.ID}, Limit: 1, Offset: 0}

		deliveryRepository.On("List", mock.Anything, listOptions).Return([]*postmand.Delivery{expectedDelivery}, nil)
		webhooks, err := webhookService.List(ctx, listOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDelivery, webhooks[0])
		deliveryRepository.AssertExpectations(t)
	})

	t.Run("Create", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		webhookService := NewDelivery(deliveryRepository)
		delivery := &postmand.Delivery{ID: uuid.New()}

		deliveryRepository.On("Create", mock.Anything, delivery).Return(nil)
		err := webhookService.Create(ctx, delivery)
		assert.Nil(t, err)
		deliveryRepository.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		webhookService := NewDelivery(deliveryRepository)
		delivery := &postmand.Delivery{ID: uuid.New()}

		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryRepository.On("Get", mock.Anything, getOptions).Return(delivery, nil)
		deliveryRepository.On("Update", mock.Anything, delivery).Return(nil)
		err := webhookService.Update(ctx, delivery)
		assert.Nil(t, err)
		deliveryRepository.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		webhookService := NewDelivery(deliveryRepository)
		delivery := &postmand.Delivery{ID: uuid.New()}

		deliveryRepository.On("Delete", mock.Anything, delivery.ID).Return(nil)
		err := webhookService.Delete(ctx, delivery.ID)
		assert.Nil(t, err)
		deliveryRepository.AssertExpectations(t)
	})
}
