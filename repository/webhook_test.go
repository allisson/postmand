package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/crypitor/postmand"
)

func makeWebhook() postmand.Webhook {
	return postmand.Webhook{
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
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("Create webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(ctx, &webhook)
		assert.Nil(t, err)
	})

	t.Run("Update webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(ctx, &webhook)
		assert.Nil(t, err)

		webhook.ValidStatusCodes = pq.Int32Array{200, 201, 204}
		err = th.webhookRepository.Update(ctx, &webhook)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		webhookFromRepository, err := th.webhookRepository.Get(ctx, options)
		assert.Nil(t, err)
		assert.Equal(t, pq.Int32Array{200, 201, 204}, webhookFromRepository.ValidStatusCodes)
	})

	t.Run("Delete webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(ctx, &webhook)
		assert.Nil(t, err)

		err = th.webhookRepository.Delete(ctx, webhook.ID)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		_, err = th.webhookRepository.Get(ctx, options)
		assert.Equal(t, postmand.ErrWebhookNotFound, err)
	})

	t.Run("Get webhook", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		err := th.webhookRepository.Create(ctx, &webhook)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		webhookFromRepository, err := th.webhookRepository.Get(ctx, options)
		assert.Nil(t, err)
		assert.Equal(t, webhook.ID, webhookFromRepository.ID)
	})

	t.Run("List webhooks", func(t *testing.T) {
		th := newTestHelper()
		defer th.db.Close()

		webhook1 := makeWebhook()
		err := th.webhookRepository.Create(ctx, &webhook1)
		assert.Nil(t, err)

		webhook2 := makeWebhook()
		err = th.webhookRepository.Create(ctx, &webhook2)
		assert.Nil(t, err)

		options := postmand.RepositoryListOptions{Limit: 1, Offset: 0, OrderBy: "created_at", Order: "DESC"}
		webhooks, err := th.webhookRepository.List(ctx, options)
		assert.Nil(t, err)
		assert.Len(t, webhooks, 1)
		assert.Equal(t, webhook2.ID, webhooks[0].ID)
	})
}
