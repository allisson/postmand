package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/postmand"
	"github.com/allisson/postmand/mocks"
)

func TestDeliveryAttempt(t *testing.T) {
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		deliveryAttemptRepository := &mocks.DeliveryAttemptRepository{}
		deliveryAttemptService := NewDeliveryAttempt(deliveryAttemptRepository)
		expectedDeliveryAttempt := &postmand.DeliveryAttempt{ID: uuid.New()}
		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": expectedDeliveryAttempt.ID}}

		deliveryAttemptRepository.On("Get", mock.Anything, getOptions).Return(expectedDeliveryAttempt, nil)
		delivery, err := deliveryAttemptService.Get(ctx, getOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveryAttempt, delivery)
		deliveryAttemptRepository.AssertExpectations(t)
	})

	t.Run("List", func(t *testing.T) {
		deliveryAttemptRepository := &mocks.DeliveryAttemptRepository{}
		deliveryAttemptService := NewDeliveryAttempt(deliveryAttemptRepository)
		expectedDeliveryAttempt := &postmand.DeliveryAttempt{ID: uuid.New()}
		listOptions := postmand.RepositoryListOptions{Filters: map[string]interface{}{"id": expectedDeliveryAttempt.ID}, Limit: 1, Offset: 0}

		deliveryAttemptRepository.On("List", mock.Anything, listOptions).Return([]*postmand.DeliveryAttempt{expectedDeliveryAttempt}, nil)
		webhooks, err := deliveryAttemptService.List(ctx, listOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedDeliveryAttempt, webhooks[0])
		deliveryAttemptRepository.AssertExpectations(t)
	})
}
