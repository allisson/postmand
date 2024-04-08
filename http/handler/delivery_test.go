package handler

import (
	"encoding/json"
	nethttp "net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/crypitor/postmand"
	"github.com/crypitor/postmand/http"
	"github.com/crypitor/postmand/mocks"
)

func makeDelivery() postmand.Delivery {
	deliveryID, _ := uuid.Parse("b919ca2c-6b0f-4a22-a61f-8c882ee69323")
	webhookID, _ := uuid.Parse("cd9b7318-36c6-4534-be84-fe78042aeaf2")

	return postmand.Delivery{
		ID:        deliveryID,
		WebhookID: webhookID,
		Payload:   `{}`,
	}
}

func TestDelivery(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("List", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		listOptions := postmand.RepositoryListOptions{Filters: map[string]interface{}{}, Limit: 50, Offset: 0, OrderBy: "created_at", Order: "desc"}
		deliveryHandler := NewDelivery(deliveryService, logger)
		router := http.NewRouter(logger)
		router.Get("/v1/deliveries", deliveryHandler.List)

		deliveryService.On("List", mock.Anything, listOptions).Return([]*postmand.Delivery{{}}, nil)
		apitest.New().
			Handler(router).
			Get("/v1/deliveries").
			Expect(t).
			Body(`{"deliveries":[{"id":"00000000-0000-0000-0000-000000000000","webhook_id":"00000000-0000-0000-0000-000000000000","payload":"","scheduled_at":"0001-01-01T00:00:00Z","delivery_attempts":0,"status":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}],"limit":50,"offset":0}`).
			Status(nethttp.StatusOK).
			End()

		deliveryService.AssertExpectations(t)
	})

	t.Run("Get", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		delivery := makeDelivery()
		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.ID}}
		deliveryHandler := NewDelivery(deliveryService, logger)
		router := http.NewRouter(logger)
		router.Get("/v1/deliveries/{delivery_id}", deliveryHandler.Get)

		deliveryService.On("Get", mock.Anything, getOptions).Return(&delivery, nil)
		apitest.New().
			Handler(router).
			Get("/v1/deliveries/b919ca2c-6b0f-4a22-a61f-8c882ee69323").
			Expect(t).
			Body(`{"created_at":"0001-01-01T00:00:00Z", "delivery_attempts":0, "id":"b919ca2c-6b0f-4a22-a61f-8c882ee69323", "payload":"{}", "scheduled_at":"0001-01-01T00:00:00Z", "status":"", "updated_at":"0001-01-01T00:00:00Z", "webhook_id":"cd9b7318-36c6-4534-be84-fe78042aeaf2"}`).
			Status(nethttp.StatusOK).
			End()

		deliveryService.AssertExpectations(t)
	})

	t.Run("Create with malformed request body", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		deliveryHandler := NewDelivery(deliveryService, logger)
		router := http.NewRouter(logger)
		router.Post("/v1/deliveries", deliveryHandler.Create)

		apitest.New().
			Handler(router).
			Post("/v1/deliveries").
			JSON(`{`).
			Expect(t).
			Body(`{"code":3, "message":"malformed request body"}`).
			Status(nethttp.StatusBadRequest).
			End()

		deliveryService.AssertExpectations(t)
	})

	t.Run("Create with valid body", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		deliveryHandler := NewDelivery(deliveryService, logger)
		delivery := makeDelivery()
		jsonDelivery, _ := json.Marshal(&delivery)
		router := http.NewRouter(logger)
		router.Post("/v1/deliveries", deliveryHandler.Create)

		deliveryService.On("Create", mock.Anything, &delivery).Return(nil)
		apitest.New().
			Handler(router).
			Post("/v1/deliveries").
			JSON(jsonDelivery).
			Expect(t).
			Body(`{"created_at":"0001-01-01T00:00:00Z", "delivery_attempts":0, "id":"b919ca2c-6b0f-4a22-a61f-8c882ee69323", "payload":"{}", "scheduled_at":"0001-01-01T00:00:00Z", "status":"", "updated_at":"0001-01-01T00:00:00Z", "webhook_id":"cd9b7318-36c6-4534-be84-fe78042aeaf2"}`).
			Status(nethttp.StatusCreated).
			End()

		deliveryService.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		deliveryService := &mocks.DeliveryService{}
		deliveryHandler := NewDelivery(deliveryService, logger)
		delivery := makeDelivery()
		router := http.NewRouter(logger)
		router.Delete("/v1/deliveries/{delivery_id}", deliveryHandler.Delete)

		deliveryService.On("Delete", mock.Anything, delivery.ID).Return(nil)
		apitest.New().
			Handler(router).
			Delete("/v1/deliveries/b919ca2c-6b0f-4a22-a61f-8c882ee69323").
			Expect(t).
			Status(nethttp.StatusNoContent).
			End()

		deliveryService.AssertExpectations(t)
	})
}
