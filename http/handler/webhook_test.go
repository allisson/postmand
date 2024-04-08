package handler

import (
	"encoding/json"
	nethttp "net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/allisson/postmand"
	"github.com/allisson/postmand/http"
	"github.com/allisson/postmand/mocks"
)

func makeWebhook() postmand.Webhook {
	webhookID, _ := uuid.Parse("cd9b7318-36c6-4534-be84-fe78042aeaf2")
	return postmand.Webhook{
		ID:                     webhookID,
		Name:                   "Test",
		URL:                    "https://httpbin.org/post",
		ContentType:            "application/json",
		ValidStatusCodes:       pq.Int32Array{200, 201},
		SecretToken:            "",
		Authorization:          "",
		Active:                 true,
		MaxDeliveryAttempts:    1,
		DeliveryAttemptTimeout: 1,
		RetryMinBackoff:        1,
		RetryMaxBackoff:        1,
	}
}

func TestWebhook(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("List", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		listOptions := postmand.RepositoryListOptions{Filters: map[string]interface{}{}, Limit: 50, Offset: 0, OrderBy: "name", Order: "asc"}
		webhookHandler := NewWebhook(webhookService, logger)
		router := http.NewRouter(logger)
		router.Get("/v1/webhooks", webhookHandler.List)

		webhookService.On("List", mock.Anything, listOptions).Return([]*postmand.Webhook{{}}, nil)
		apitest.New().
			Handler(router).
			Get("/v1/webhooks").
			Expect(t).
			Body(`{"webhooks":[{"id":"00000000-0000-0000-0000-000000000000","name":"","url":"","content_type":"","valid_status_codes":null,"secret_token":"","active":false,"max_delivery_attempts":0,"delivery_attempt_timeout":0,"retry_min_backoff":0,"retry_max_backoff":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}],"limit":50,"offset":0}`).
			Status(nethttp.StatusOK).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Get", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhook := makeWebhook()
		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		webhookHandler := NewWebhook(webhookService, logger)
		router := http.NewRouter(logger)
		router.Get("/v1/webhooks/{webhook_id}", webhookHandler.Get)

		webhookService.On("Get", mock.Anything, getOptions).Return(&webhook, nil)
		apitest.New().
			Handler(router).
			Get("/v1/webhooks/cd9b7318-36c6-4534-be84-fe78042aeaf2").
			Expect(t).
			Body(`{"active":true, "content_type":"application/json", "created_at":"0001-01-01T00:00:00Z", "delivery_attempt_timeout":1, "id":"cd9b7318-36c6-4534-be84-fe78042aeaf2", "max_delivery_attempts":1, "name":"Test", "retry_max_backoff":1, "retry_min_backoff":1, "secret_token":"", "updated_at":"0001-01-01T00:00:00Z", "url":"https://httpbin.org/post", "valid_status_codes":[200, 201]}`).
			Status(nethttp.StatusOK).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Create with malformed request body", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		router := http.NewRouter(logger)
		router.Post("/v1/webhooks", webhookHandler.Create)

		apitest.New().
			Handler(router).
			Post("/v1/webhooks").
			JSON(`{`).
			Expect(t).
			Body(`{"code":3, "message":"malformed request body"}`).
			Status(nethttp.StatusBadRequest).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Create with invalid body", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		router := http.NewRouter(logger)
		router.Post("/v1/webhooks", webhookHandler.Create)

		apitest.New().
			Handler(router).
			Post("/v1/webhooks").
			JSON(`{}`).
			Expect(t).
			Body(`{"code":4, "details":"content_type: cannot be blank; delivery_attempt_timeout: cannot be blank; max_delivery_attempts: cannot be blank; name: cannot be blank; retry_max_backoff: cannot be blank; retry_min_backoff: cannot be blank; url: cannot be blank; valid_status_codes: cannot be blank.", "message":"request validation failed"}`).
			Status(nethttp.StatusBadRequest).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Create with valid body", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		webhook := makeWebhook()
		jsonWebhook, _ := json.Marshal(&webhook)
		router := http.NewRouter(logger)
		router.Post("/v1/webhooks", webhookHandler.Create)

		webhookService.On("Create", mock.Anything, &webhook).Return(nil)
		apitest.New().
			Handler(router).
			Post("/v1/webhooks").
			JSON(jsonWebhook).
			Expect(t).
			Body(`{"active":true, "content_type":"application/json", "created_at":"0001-01-01T00:00:00Z", "delivery_attempt_timeout":1, "id":"cd9b7318-36c6-4534-be84-fe78042aeaf2", "max_delivery_attempts":1, "name":"Test", "retry_max_backoff":1, "retry_min_backoff":1, "secret_token":"", "updated_at":"0001-01-01T00:00:00Z", "url":"https://httpbin.org/post", "valid_status_codes": [200, 201]}`).
			Status(nethttp.StatusCreated).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Update with malformed request body", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		router := http.NewRouter(logger)
		router.Put("/v1/webhooks/{webhook_id}", webhookHandler.Update)

		apitest.New().
			Handler(router).
			Put("/v1/webhooks/cd9b7318-36c6-4534-be84-fe78042aeaf2").
			JSON(`{`).
			Expect(t).
			Body(`{"code":3, "message":"malformed request body"}`).
			Status(nethttp.StatusBadRequest).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Update with invalid body", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		router := http.NewRouter(logger)
		router.Put("/v1/webhooks/{webhook_id}", webhookHandler.Update)

		apitest.New().
			Handler(router).
			Put("/v1/webhooks/cd9b7318-36c6-4534-be84-fe78042aeaf2").
			JSON(`{}`).
			Expect(t).
			Body(`{"code":4, "details":"content_type: cannot be blank; delivery_attempt_timeout: cannot be blank; max_delivery_attempts: cannot be blank; name: cannot be blank; retry_max_backoff: cannot be blank; retry_min_backoff: cannot be blank; url: cannot be blank; valid_status_codes: cannot be blank.", "message":"request validation failed"}`).
			Status(nethttp.StatusBadRequest).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Update with valid body", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		webhook := makeWebhook()
		jsonWebhook, _ := json.Marshal(&webhook)
		router := http.NewRouter(logger)
		router.Put("/v1/webhooks/{webhook_id}", webhookHandler.Update)

		webhookService.On("Update", mock.Anything, &webhook).Return(nil)
		apitest.New().
			Handler(router).
			Put("/v1/webhooks/cd9b7318-36c6-4534-be84-fe78042aeaf2").
			JSON(jsonWebhook).
			Expect(t).
			Body(`{"active":true, "content_type":"application/json", "created_at":"0001-01-01T00:00:00Z", "delivery_attempt_timeout":1, "id":"cd9b7318-36c6-4534-be84-fe78042aeaf2", "max_delivery_attempts":1, "name":"Test", "retry_max_backoff":1, "retry_min_backoff":1, "secret_token":"", "updated_at":"0001-01-01T00:00:00Z", "url":"https://httpbin.org/post", "valid_status_codes":[200, 201]}`).
			Status(nethttp.StatusOK).
			End()

		webhookService.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		webhookService := &mocks.WebhookService{}
		webhookHandler := NewWebhook(webhookService, logger)
		webhook := makeWebhook()
		router := http.NewRouter(logger)
		router.Delete("/v1/webhooks/{webhook_id}", webhookHandler.Delete)

		webhookService.On("Delete", mock.Anything, webhook.ID).Return(nil)
		apitest.New().
			Handler(router).
			Delete("/v1/webhooks/cd9b7318-36c6-4534-be84-fe78042aeaf2").
			Expect(t).
			Status(nethttp.StatusNoContent).
			End()

		webhookService.AssertExpectations(t)
	})
}
