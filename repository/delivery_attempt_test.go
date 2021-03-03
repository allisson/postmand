package repository

import (
	"testing"

	"github.com/allisson/postmand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func makeDeliveryAttempt() *postmand.DeliveryAttempt {
	return &postmand.DeliveryAttempt{
		ID:                 uuid.New(),
		ResponseStatusCode: 201,
		ExecutionDuration:  150,
		Success:            true,
		Error:              "",
	}
}

func TestDeliveryAttempt(t *testing.T) {
	t.Run("Create delivery attempt", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		deliveryAttempt := makeDeliveryAttempt()
		deliveryAttempt.WebhookID = webhook.ID
		deliveryAttempt.DeliveryID = delivery.ID
		err = th.deliveryAttemptRepository.Create(deliveryAttempt)
		assert.Nil(t, err)
	})

	t.Run("Get delivery attempt", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		deliveryAttempt := makeDeliveryAttempt()
		deliveryAttempt.WebhookID = webhook.ID
		deliveryAttempt.DeliveryID = delivery.ID
		err = th.deliveryAttemptRepository.Create(deliveryAttempt)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": deliveryAttempt.ID}}
		deliveryAttemptFromRepository, err := th.deliveryAttemptRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, deliveryAttempt.ID, deliveryAttemptFromRepository.ID)
	})

	t.Run("List delivery attempts", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		deliveryAttempt1 := makeDeliveryAttempt()
		deliveryAttempt1.WebhookID = webhook.ID
		deliveryAttempt1.DeliveryID = delivery.ID
		err = th.deliveryAttemptRepository.Create(deliveryAttempt1)
		assert.Nil(t, err)

		deliveryAttempt2 := makeDeliveryAttempt()
		deliveryAttempt2.WebhookID = webhook.ID
		deliveryAttempt2.DeliveryID = delivery.ID
		err = th.deliveryAttemptRepository.Create(deliveryAttempt2)
		assert.Nil(t, err)

		options := postmand.RepositoryListOptions{Limit: 1, Offset: 1, OrderBy: "created_at DESC"}
		deliveryAttempts, err := th.deliveryAttemptRepository.List(options)
		assert.Nil(t, err)
		assert.Len(t, deliveryAttempts, 1)
		assert.Equal(t, deliveryAttempt2.ID, deliveryAttempts[0].ID)
	})
}
