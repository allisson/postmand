package repository

import (
	"database/sql"
	"testing"

	"github.com/allisson/postmand"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func makeWebhook() *postmand.Webhook {
	return &postmand.Webhook{
		ID:                     uuid.New(),
		Name:                   "Test",
		URL:                    "https://httpbin.org/post",
		ContentType:            "application/json",
		Active:                 true,
		ValidStatusCodes:       pq.Int32Array{200, 201},
		MaxDeliveryAttempts:    1,
		DeliveryAttemptTimeout: 1,
		RetryMinBackoff:        1,
		RetryMaxBackoff:        1,
	}
}

func TestTransaction(t *testing.T) {
	t.Run("Create webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)
	})

	t.Run("Update webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		webhook.ValidStatusCodes = pq.Int32Array{200, 201, 204}
		err = th.webhookRepository.Update(webhook)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		webhookFromRepository, err := th.webhookRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, pq.Int32Array{200, 201, 204}, webhookFromRepository.ValidStatusCodes)
	})

	t.Run("Delete webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		err = th.webhookRepository.Delete(webhook.ID)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		_, err = th.webhookRepository.Get(options)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Get webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		webhookFromRepository, err := th.webhookRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, webhook.ID, webhookFromRepository.ID)
	})

	t.Run("List webhooks", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook1 := makeWebhook()
		err := th.webhookRepository.Create(webhook1)
		assert.Nil(t, err)

		webhook2 := makeWebhook()
		err = th.webhookRepository.Create(webhook2)
		assert.Nil(t, err)

		options := postmand.RepositoryListOptions{Limit: 1, Offset: 1, OrderBy: "created_at DESC"}
		webhooks, err := th.webhookRepository.List(options)
		assert.Nil(t, err)
		assert.Len(t, webhooks, 1)
		assert.Equal(t, webhook2.ID, webhooks[0].ID)
	})
}
