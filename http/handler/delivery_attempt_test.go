package handler

import (
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

func makeDeliveryAttempt() postmand.DeliveryAttempt {
	deliveryAttemptID, _ := uuid.Parse("97087247-d89d-410e-b915-740b4c6d9d99")
	deliveryID, _ := uuid.Parse("b919ca2c-6b0f-4a22-a61f-8c882ee69323")
	webhookID, _ := uuid.Parse("cd9b7318-36c6-4534-be84-fe78042aeaf2")

	return postmand.DeliveryAttempt{
		ID:         deliveryAttemptID,
		DeliveryID: deliveryID,
		WebhookID:  webhookID,
	}
}

func TestDeliveryAttempt(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("List", func(t *testing.T) {
		deliveryAttemptService := &mocks.DeliveryAttemptService{}
		listOptions := postmand.RepositoryListOptions{Filters: map[string]interface{}{}, Limit: 50, Offset: 0, OrderBy: "created_at", Order: "desc"}
		deliveryAttemptHandler := NewDeliveryAttempt(deliveryAttemptService, logger)
		router := http.NewRouter(logger)
		router.Get("/v1/delivery-attempts", deliveryAttemptHandler.List)

		deliveryAttemptService.On("List", mock.Anything, listOptions).Return([]*postmand.DeliveryAttempt{{}}, nil)
		apitest.New().
			Handler(router).
			Get("/v1/delivery-attempts").
			Expect(t).
			Body(`{"delivery_attempts":[{"id":"00000000-0000-0000-0000-000000000000","webhook_id":"00000000-0000-0000-0000-000000000000","delivery_id":"00000000-0000-0000-0000-000000000000","raw_request":"", "raw_response":"","response_status_code":0,"execution_duration":0,"success":false,"error":"","created_at":"0001-01-01T00:00:00Z"}],"limit":50,"offset":0}`).
			Status(nethttp.StatusOK).
			End()

		deliveryAttemptService.AssertExpectations(t)
	})

	t.Run("Get", func(t *testing.T) {
		deliveryAttemptService := &mocks.DeliveryAttemptService{}
		deliveryAttempt := makeDeliveryAttempt()
		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": deliveryAttempt.ID}}
		deliveryAttemptHandler := NewDeliveryAttempt(deliveryAttemptService, logger)
		router := http.NewRouter(logger)
		router.Get("/v1/delivery-attempts/{delivery_attempt_id}", deliveryAttemptHandler.Get)

		deliveryAttemptService.On("Get", mock.Anything, getOptions).Return(&deliveryAttempt, nil)
		apitest.New().
			Handler(router).
			Get("/v1/delivery-attempts/97087247-d89d-410e-b915-740b4c6d9d99").
			Expect(t).
			Body(`{"id":"97087247-d89d-410e-b915-740b4c6d9d99","webhook_id":"cd9b7318-36c6-4534-be84-fe78042aeaf2","delivery_id":"b919ca2c-6b0f-4a22-a61f-8c882ee69323","raw_request":"", "raw_response":"","response_status_code":0,"execution_duration":0,"success":false,"error":"","created_at":"0001-01-01T00:00:00Z"}`).
			Status(nethttp.StatusOK).
			End()

		deliveryAttemptService.AssertExpectations(t)
	})
}
