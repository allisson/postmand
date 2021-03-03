package repository

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
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
		Status:           postmand.DeliveryStatusPending,
	}
}

func TestDispatchToURL(t *testing.T) {
	t.Run("Invalid webhook url", func(t *testing.T) {
		webhook := makeWebhook()
		webhook.URL = "http://localhost:9999"
		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID

		dr := dispatchToURL(webhook, delivery)
		assert.False(t, dr.Success)
		assert.Equal(t, `Post "http://localhost:9999": dial tcp [::1]:9999: connect: connection refused`, dr.Error)
	})

	t.Run("Invalid response status code", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		webhook := makeWebhook()
		webhook.URL = httpServer.URL
		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID

		dr := dispatchToURL(webhook, delivery)
		assert.NotEqual(t, "", dr.RawResponse)
		assert.Equal(t, http.StatusNoContent, dr.ResponseStatusCode)
		assert.False(t, dr.Success)
		assert.Equal(t, "", dr.Error)
	})

	t.Run("Valid response status code", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		webhook := makeWebhook()
		webhook.URL = httpServer.URL
		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID

		dr := dispatchToURL(webhook, delivery)
		assert.NotEqual(t, "", dr.RawResponse)
		assert.Equal(t, http.StatusOK, dr.ResponseStatusCode)
		assert.True(t, dr.Success)
		assert.Equal(t, "", dr.Error)
	})
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

		delivery.Status = postmand.DeliveryStatusPending
		err = th.deliveryRepository.Update(delivery)
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryFromRepository, err := th.deliveryRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, postmand.DeliveryStatusPending, deliveryFromRepository.Status)
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

	t.Run("Dispatch delivery succeeded", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		webhook.URL = httpServer.URL
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		err = th.deliveryRepository.Dispatch()
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryFromRepository, err := th.deliveryRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, 1, deliveryFromRepository.DeliveryAttempts)
		assert.Equal(t, postmand.DeliveryStatusSucceeded, deliveryFromRepository.Status)

		options = postmand.RepositoryGetOptions{Filters: map[string]interface{}{"delivery_id": delivery.ID}}
		deliveryAttemptFromRepository, err := th.deliveryAttemptRepository.Get(options)
		assert.Nil(t, err)
		assert.True(t, deliveryAttemptFromRepository.Success)
	})

	t.Run("Dispatch delivery retry", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		webhook.MaxDeliveryAttempts = 2
		webhook.URL = httpServer.URL
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		err = th.deliveryRepository.Dispatch()
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryFromRepository, err := th.deliveryRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, 1, deliveryFromRepository.DeliveryAttempts)
		assert.Equal(t, postmand.DeliveryStatusPending, deliveryFromRepository.Status)
		assert.True(t, deliveryFromRepository.ScheduledAt.After(delivery.ScheduledAt))

		options = postmand.RepositoryGetOptions{Filters: map[string]interface{}{"delivery_id": delivery.ID}}
		deliveryAttemptFromRepository, err := th.deliveryAttemptRepository.Get(options)
		assert.Nil(t, err)
		assert.False(t, deliveryAttemptFromRepository.Success)
	})

	t.Run("Dispatch delivery failed", func(t *testing.T) {
		httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
			// nolint:errcheck
			w.Write([]byte("OK"))
		}))
		defer httpServer.Close()

		th := newTestHelper()
		defer th.db.Close()

		webhook := makeWebhook()
		webhook.URL = httpServer.URL
		err := th.webhookRepository.Create(webhook)
		assert.Nil(t, err)

		delivery := makeDelivery()
		delivery.WebhookID = webhook.ID
		err = th.deliveryRepository.Create(delivery)
		assert.Nil(t, err)

		err = th.deliveryRepository.Dispatch()
		assert.Nil(t, err)

		options := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryFromRepository, err := th.deliveryRepository.Get(options)
		assert.Nil(t, err)
		assert.Equal(t, 1, deliveryFromRepository.DeliveryAttempts)
		assert.Equal(t, postmand.DeliveryStatusFailed, deliveryFromRepository.Status)

		options = postmand.RepositoryGetOptions{Filters: map[string]interface{}{"delivery_id": delivery.ID}}
		deliveryAttemptFromRepository, err := th.deliveryAttemptRepository.Get(options)
		assert.Nil(t, err)
		assert.False(t, deliveryAttemptFromRepository.Success)
	})
}
