package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/allisson/postmand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func makeDelivery() *postmand.Delivery {
	return &postmand.Delivery{
		ID:               uuid.New(),
		Payload:          `{"success": true}`,
		ScheduledAt:      time.Now().UTC(),
		DeliveryAttempts: 0,
		Status:           postmand.DeliveryStatusTodo,
	}
}

func TestDelivery(t *testing.T) {
	t.Run("Create delivery", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)
	})

	t.Run("Update delivery", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		delivery.Status = postmand.DeliveryStatusDoing
		err = th.deliveryRepository.Update(delivery)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryFromRepository, err := th.deliveryRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, postmand.DeliveryStatusDoing, deliveryFromRepository.Status)
	})

	t.Run("Delete delivery", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		err = th.deliveryRepository.Delete(delivery.ID)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		_, err = th.deliveryRepository.Get(options)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Get delivery", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryFromRepository, err := th.deliveryRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, delivery.ID, deliveryFromRepository.ID)
	})

	t.Run("List deliveries", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery1 := makeDelivery()
		delivery1.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery1)
		assert.Nil(t, err)

		delivery2 := makeDelivery()
		delivery2.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery2)
		assert.Nil(t, err)

		options := postmand.RepositoryListOptions{Limit: 1, Offset: 1, OrderBy: "created_at DESC"}
		deliveries, err := th.deliveryRepository.List(options)
		assert.Nil(t, err)
		assert.Len(t, deliveries, 1)
		assert.Equal(t, delivery2.ID, deliveries[0].ID)
	})
}
